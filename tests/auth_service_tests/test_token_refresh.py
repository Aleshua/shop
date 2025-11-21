
def test_token_refresh(auth_client, logged_user):    
    resp = auth_client.refresh_token(
        logged_user["session_data"]["refresh_token"]
    )
    
    assert resp.status_code == 200
    
    body = resp.json()
    assert isinstance(body, dict)
    assert isinstance(body["data"], dict)
    assert isinstance(body["data"]["access_token"], str)

def test_fake_token_refresh(auth_client):    
    resp = auth_client.refresh_token(
        "lol"
    )
    
    assert resp.status_code == 404
    
def test_token_in_blacklist_refresh(auth_client, logged_user):
    
    resp = auth_client.logout(
        logged_user["session_data"]["refresh_token"]
    )
    
    resp = auth_client.refresh_token(
        logged_user["session_data"]["refresh_token"]
    )
    
    assert resp.status_code == 401
