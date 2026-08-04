package email_store

import "sync"

type EmailStoreInfo struct {
	Email string
	Code  string
}

var emailVerifyStore = sync.Map{}

func Set(id, email, code string) {
	emailVerifyStore.Store(id, EmailStoreInfo{
		Email: email,
		Code:  code,
	})
}

func Verify(id, code string) (info EmailStoreInfo, ok bool) {
	value, ok := emailVerifyStore.Load(id)
	emailVerifyStore.Delete(id)
	if !ok {
		return
	}
	info, ok = value.(EmailStoreInfo)
	if !ok {
		return
	}
	if info.Code != code {
		ok = false
		return
	}
	return
}
