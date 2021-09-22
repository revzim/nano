package drivers

import (
	"testing"
	"time"
)

func TestFirebase(t *testing.T) {
	const credentialsPath = ""
	const email = ""
	const password = ""
	const displayName = ""
	azfb, err := NewFirebaseApp("")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Logf("initializing new firebase user: %s | %s\n", email, displayName)
	newUser, err := azfb.InitNewUser(email, password, displayName)
	if err != nil {
		t.Log("new user init failed!", err)
		t.Fail()
	}
	t.Log(newUser.UID)
	userData, err := azfb.GetUser(email)
	if err != nil {
		t.Log("get user error:", err)
		t.Fail()
	}
	t.Logf("initialized user with uid: %s\n", userData.UID)

	var newTokenStr string
	nowTime := time.Now().Unix()
	const duration = 10
	customData := map[string]interface{}{
		"id":         userData.UID,
		"created_at": nowTime,
		"expiry":     nowTime + duration,
	}
	newTokenStr, err = azfb.InitCustomToken(userData.UID, customData)
	if err != nil {
		t.Log("custom token gen error:", err)
		t.Fail()
	}
	t.Logf("new custom token initialized: %s\n", newTokenStr)

	// var newUserData *auth.UserRecord
	newUserData, err := azfb.UpdateUser(userData.UID, FirebaseUpdateOpts{
		Email: "updatedEmail@test.com",
		// Email: fmt.Sprintf("%s@test.com", uuid.New().String()),
	})
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Logf("email updated from %s to %s\n", userData.Email, newUserData.Email)

	err = azfb.RevokeRefreshTokens(newUserData.UID)
	if err != nil {
		t.Log("token revokation error:", err)
		t.Fail()
	}
	t.Logf("%s's user token was succesfully revoked!\n", newUserData.UID)

	t.Logf("deleting user: %s ...\n", newUserData.UID)
	err = azfb.DeleteUser(newUserData.UID)
	if err != nil {
		t.Log("user deletion error:", err)
		t.Fail()
	}
	t.Logf("successfully deleted user: %s\n", newUserData.UID)
}
