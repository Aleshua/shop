
def test_login(auth_client, registered_user):
    resp = auth_client.login(
        registered_user["email"],
        registered_user["password"],
    )
    assert resp.status_code == 200
    
    body = resp.json()
    assert isinstance(body, dict)
    assert isinstance(body["data"], dict)
    assert isinstance(body["data"]["refresh_token"], str)
    assert isinstance(body["data"]["access_token"], str)

def test_login_with_wrong_password(auth_client, registered_user):
    resp = auth_client.login(
        registered_user["email"],
        "lol",
    )
    assert resp.status_code == 401

def test_login_with_wrong_email(auth_client, registered_user):
    resp = auth_client.login(
        "lol",
        registered_user["password"],
    )
    assert resp.status_code == 400
