import re
import time
import quopri
import requests
from typing import Dict, Any, List

def get_latest_confirmation_code(recipient_email: str) -> str:
    MAILHOG_API_URL = "http://localhost:8025/api/v2/messages"
    ATTEMPT_INTERVAL = 0.1
    MAX_ATTEMPTS = int(2 / ATTEMPT_INTERVAL)
    
    for _ in range(MAX_ATTEMPTS):
        time.sleep(ATTEMPT_INTERVAL)
    
        resp: requests.Response = requests.get(MAILHOG_API_URL)
        resp.raise_for_status()
        
        messages: List[Dict[str, Any]] = resp.json().get("items", [])
        
        for msg in messages:
            recipients_list = msg.get("To", [])
            
            found_recipient = False
            for recipient in recipients_list:
                full_address = f"{recipient.get('Mailbox', '')}@{recipient.get('Domain', '')}"
                if recipient_email == full_address:
                    found_recipient = True
                    break
                    
            if found_recipient:
                
                body_bytes = msg["Content"]["Body"].encode('latin-1')
                decoded_body = quopri.decodestring(body_bytes).decode('utf-8')

                match = re.search(r'Ваш код подтверждения:\s*(\w+)', decoded_body)
                
                if match:
                    return match.group(1)
                else:
                    raise ValueError(
                        f"не удалось найти шаблон кода в теле письма для {recipient_email}. тело: {decoded_body[:100]}",
                    )

    raise TimeoutError(f"письмо для {recipient_email} не найдено в Mailhog")