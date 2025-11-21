from tests.utils.mailhog import get_latest_confirmation_code

def test_register(auth_client, random_user_data):
    resp = auth_client.register(
        random_user_data["email"], 
        random_user_data["password"]
    )
    
    assert resp.status_code == 200
    
    body = resp.json()
    
    assert isinstance(body, dict)
    assert isinstance(body["data"], dict)
    assert isinstance(body["data"]["id"], int)

def test_register_repeat(auth_client, random_user_data):
    resp = auth_client.register(
        random_user_data["email"], 
        random_user_data["password"]
    )
    
    resp = auth_client.register(
        random_user_data["email"], 
        random_user_data["password"]
    )
    
    assert resp.status_code == 200

def test_register_confirm_email(auth_client, random_user_data):
    resp = auth_client.register(
        random_user_data["email"], 
        random_user_data["password"]
    )
    
    resp = auth_client.confirm_email(
        resp.json()["data"]["id"],
        get_latest_confirmation_code(random_user_data["email"]),
    )    

    assert resp.status_code == 200

def test_register_email_conflict(auth_client, registered_user):        
    resp = auth_client.register(
        registered_user["email"], 
        registered_user["password"]
    )
    
    assert resp.status_code == 409

def test_register_confirm_email_with_wrong_code(auth_client, random_user_data):
    resp = auth_client.register(
        random_user_data["email"], 
        random_user_data["password"]
    )
    
    resp = auth_client.confirm_email(
        resp.json()["data"]["id"],
        "lol",
    )  
    
    assert resp.status_code == 400

def test_register_confirm_email_when_attempts_exceeded(auth_client, random_user_data):
    resp = auth_client.register(
        random_user_data["email"], 
        random_user_data["password"]
    )
    
    user_id = resp.json()["data"]["id"]
    
    for _ in range(3):
        resp = auth_client.confirm_email(
            user_id,
            "lol",
        )  
        assert resp.status_code == 400
    
    resp = auth_client.confirm_email(
        user_id,
        "lol",
    )  
    assert resp.status_code == 403

def test_register_resend_code(auth_client, random_user_data):
    resp = auth_client.register(
        random_user_data["email"], 
        random_user_data["password"]
    )
    
    user_id = resp.json()["data"]["id"]
    
    for _ in range(3):
        resp = auth_client.confirm_email(
            user_id,
            "lol",
        )  
        assert resp.status_code == 400
    
    resp = auth_client.resend_code(
        user_id,
    ) 
    assert resp.status_code == 200
    
    resp = auth_client.confirm_email(
        user_id,
        get_latest_confirmation_code(random_user_data["email"]),
    )  
    assert resp.status_code == 200
