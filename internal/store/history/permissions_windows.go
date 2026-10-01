//go:build windows

package history

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

func isUnsafePathLink(path string, info os.FileInfo) (bool, error) {
	if info.Mode()&os.ModeSymlink != 0 {
		return true, nil
	}
	name, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return false, err
	}
	attributes, err := windows.GetFileAttributes(name)
	if err != nil {
		return false, err
	}
	return attributes&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0, nil
}

func secureDirectory(path string, writable, newlyCreated bool) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return errors.New("history directory is not a real directory")
	}
	if writable {
		if err := setPrivateACL(path, true, newlyCreated); err != nil {
			return err
		}
	}
	return verifyPrivateACL(path, true)
}

func secureFile(path string, writable, newlyCreated bool) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return errors.New("history database is not a regular file")
	}
	if writable {
		if err := setPrivateACL(path, false, newlyCreated); err != nil {
			return err
		}
	}
	return verifyPrivateACL(path, false)
}

func secureSQLiteSidecars(database string, allowNewOwner bool) error {
	for _, suffix := range []string{"-wal", "-shm"} {
		path := database + suffix
		if err := secureFile(path, true, allowNewOwner); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}

func currentUserSID() (*windows.SID, error) {
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		return nil, err
	}
	return user.User.Sid, nil
}

func setPrivateACL(path string, directory, allowOwnerChange bool) error {
	currentSID, err := currentUserSID()
	if err != nil {
		return err
	}
	security, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.OWNER_SECURITY_INFORMATION)
	if err != nil {
		return err
	}
	owner, _, err := security.Owner()
	if err != nil || owner == nil {
		return errors.New("history path owner could not be verified")
	}
	ownerToSet := (*windows.SID)(nil)
	securityInfo := windows.SECURITY_INFORMATION(windows.DACL_SECURITY_INFORMATION | windows.PROTECTED_DACL_SECURITY_INFORMATION)
	if !owner.Equals(currentSID) {
		if !allowOwnerChange {
			return errors.New("history path is not owned by the current user")
		}
		ownerToSet = currentSID
		securityInfo |= windows.OWNER_SECURITY_INFORMATION
	}
	var pinner runtime.Pinner
	pinner.Pin(currentSID)
	defer pinner.Unpin()
	inheritance := uint32(0)
	if directory {
		inheritance = windows.SUB_CONTAINERS_AND_OBJECTS_INHERIT
	}
	entry := windows.EXPLICIT_ACCESS{
		AccessPermissions: windows.GENERIC_ALL,
		AccessMode:        windows.SET_ACCESS,
		Inheritance:       inheritance,
		Trustee: windows.TRUSTEE{
			TrusteeForm:  windows.TRUSTEE_IS_SID,
			TrusteeType:  windows.TRUSTEE_IS_USER,
			TrusteeValue: windows.TrusteeValueFromSID(currentSID),
		},
	}
	acl, err := windows.ACLFromEntries([]windows.EXPLICIT_ACCESS{entry}, nil)
	if err != nil {
		return err
	}
	if err := windows.SetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, securityInfo,
		ownerToSet, nil, acl, nil); err != nil {
		return err
	}
	return verifyPrivateACL(path, directory)
}

func verifyPrivateACL(path string, directory bool) error {
	currentSID, err := currentUserSID()
	if err != nil {
		return err
	}
	security, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT,
		windows.OWNER_SECURITY_INFORMATION|windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		return err
	}
	owner, _, err := security.Owner()
	if err != nil || owner == nil || !owner.Equals(currentSID) {
		return errors.New("history path is not owned by the current user")
	}
	acl, _, err := security.DACL()
	if err != nil || acl == nil {
		return errors.New("history path DACL is unavailable")
	}
	if acl.AceCount != 1 {
		entries := make([]string, 0, acl.AceCount)
		for index := uint32(0); index < uint32(acl.AceCount); index++ {
			var entry *windows.ACCESS_ALLOWED_ACE
			if err := windows.GetAce(acl, index, &entry); err != nil || entry == nil {
				entries = append(entries, fmt.Sprintf("index=%d unreadable", index))
				continue
			}
			currentUser := false
			if entry.Header.AceType == windows.ACCESS_ALLOWED_ACE_TYPE {
				entrySID := (*windows.SID)(unsafe.Pointer(&entry.SidStart))
				currentUser = entrySID.Equals(currentSID)
			}
			entries = append(entries, fmt.Sprintf("type=%d mask=%#x flags=%#x currentUser=%t", entry.Header.AceType, entry.Mask, entry.Header.AceFlags, currentUser))
		}
		return fmt.Errorf("history path DACL has %d entries %v", acl.AceCount, entries)
	}
	var ace *windows.ACCESS_ALLOWED_ACE
	if err := windows.GetAce(acl, 0, &ace); err != nil || ace == nil || ace.Header.AceType != windows.ACCESS_ALLOWED_ACE_TYPE {
		return errors.New("history path DACL does not have a simple allow entry")
	}
	wantFlags := byte(0)
	if directory {
		wantFlags = windows.OBJECT_INHERIT_ACE | windows.CONTAINER_INHERIT_ACE
	}
	if ace.Mask != windows.GENERIC_ALL {
		return errors.New("history path DACL access mask differs")
	}
	if ace.Header.AceFlags != wantFlags {
		return errors.New("history path DACL inheritance flags differ")
	}
	aceSID := (*windows.SID)(unsafe.Pointer(&ace.SidStart))
	if !aceSID.Equals(currentSID) {
		return errors.New("history path DACL principal differs")
	}
	return nil
}
