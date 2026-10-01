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
	control, _, err := security.Control()
	if err != nil || control&windows.SE_DACL_PROTECTED == 0 {
		return errors.New("history path DACL is not protected from inheritance")
	}
	acl, _, err := security.DACL()
	if err != nil || acl == nil {
		return errors.New("history path DACL is unavailable")
	}
	expectedEntries := 1
	inheritedFlags := byte(0)
	if directory {
		expectedEntries = 2
		inheritedFlags = windows.OBJECT_INHERIT_ACE | windows.CONTAINER_INHERIT_ACE | windows.INHERIT_ONLY_ACE
	}
	if int(acl.AceCount) != expectedEntries {
		return fmt.Errorf("history path DACL has %d entries", acl.AceCount)
	}
	seenCurrent, seenInherited := false, false
	for index := uint32(0); index < uint32(acl.AceCount); index++ {
		var ace *windows.ACCESS_ALLOWED_ACE
		if err := windows.GetAce(acl, index, &ace); err != nil || ace == nil || ace.Header.AceType != windows.ACCESS_ALLOWED_ACE_TYPE {
			return errors.New("history path DACL does not have a simple allow entry")
		}
		if ace.Mask != windows.GENERIC_ALL && ace.Mask != windows.ACCESS_MASK(windowsFileAllAccess) {
			return errors.New("history path DACL access mask differs")
		}
		flags := ace.Header.AceFlags
		if flags == 0 {
			if seenCurrent {
				return errors.New("history path DACL has duplicate current-object entries")
			}
			seenCurrent = true
		} else if directory && flags == inheritedFlags {
			if seenInherited {
				return errors.New("history path DACL has duplicate inherited entries")
			}
			seenInherited = true
		} else {
			return errors.New("history path DACL inheritance flags differ")
		}
		aceSID := (*windows.SID)(unsafe.Pointer(&ace.SidStart))
		if !aceSID.Equals(currentSID) {
			return errors.New("history path DACL principal differs")
		}
	}
	if !seenCurrent || directory && !seenInherited {
		return errors.New("history path DACL does not cover its storage path")
	}
	return nil
}

const windowsFileAllAccess uint32 = 0x001f01ff
