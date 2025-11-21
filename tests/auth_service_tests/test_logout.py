
def test_logout(auth_client, logged_user):    
    resp = auth_client.logout(
        logged_user["session_data"]["refresh_token"]
    )
    
    assert resp.status_code == 200

def test_logout_with_wrong_token(auth_client, logged_user):
    resp = auth_client.logout(
        "lol"
    )
    
    assert resp.status_code == 404
