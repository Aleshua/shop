import requests


def test_ping_server():
    resp = requests.get("http://localhost/auth/ping")
    
    assert resp.status_code == 200