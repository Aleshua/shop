import pytest
import psycopg2

@pytest.fixture(scope="session")
def connect_auth_database():
    conn = psycopg2.connect(
        host="localhost",
        dbname="auth",
        user="user",
        password="password"
    )
    yield conn    
    conn.close()

@pytest.fixture(autouse=True)
def cleanup(connect_auth_database):
    cur = connect_auth_database.cursor()
    cur.execute("TRUNCATE users, refresh_tokens, confirm_cods RESTART IDENTITY CASCADE;")
    connect_auth_database.commit()
    yield
