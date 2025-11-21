from .base_client import BaseClient

class AuthClient(BaseClient):

    def register(self, email: str, password: str):
        return self.post("/auth/register", json={
            "email": email,
            "password": password,
        })

    def confirm_email(self, user_id: int, code: str):
        return self.post("/auth/email/confirm", json={
            "user_id": user_id,
            "code": code,
        })
    
    def resend_code(self, user_id: int):
        return self.post("/auth/email/resend", json={
            "user_id": user_id,
        })
        
    def login(self, email: str, password: str):
        return self.post("/auth/login", json={
            "email": email,
            "password": password,
        })
    
    def logout(self, refresh_token: str):
        return self.post("/auth/logout", json={
            "refresh_token": refresh_token,
        })
    
    def refresh_token(self, refresh_token: str):
        return self.post("/auth/token/refresh", json={
            "refresh_token": refresh_token,
        })
