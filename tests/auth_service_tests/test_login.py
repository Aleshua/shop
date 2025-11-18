import requests

def user_register():
    resp = requests.post(
        url="http://localhost/auth/register",
        json={
            "email": "example@example.com",
            "password": "strong_pass_123",
        },
    )
    resp = requests.post(
        url="http://localhost/auth/email/confirm",
        json={
            "user_id": resp.json()["data"]["id"],
            "code": "0000000000",
        },
    )
    

def test_login():
    user_register()
    resp = requests.post(
        url="http://localhost/auth/login",
        json={
            "email": "example@example.com",
            "password": "strong_pass_123",
        }
    )
    assert resp.status_code == 200
    
    body = resp.json()
    assert isinstance(body, dict)
    assert isinstance(body["data"], dict)
    assert isinstance(body["data"]["refresh_token"], str)
    assert isinstance(body["data"]["access_token"], str)

def test_login_with_wrong_password():
    user_register()
    resp = requests.post(
        url="http://localhost/auth/login",
        json={
            "email": "example@example.com",
            "password": "strong_pass_124",
        }
    )
    assert resp.status_code == 401

def test_login_with_wrong_email():
    user_register()
    resp = requests.post(
        url="http://localhost/auth/login",
        json={
            "email": "example2@example.com",
            "password": "strong_pass_123",
        }
    )
    assert resp.status_code == 404
