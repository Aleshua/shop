# import psycopg2
# from typing import Any, List, Tuple

# def database_fetch_all(database: psycopg2.extensions.connection, table: str, **kwargs: Any) -> List[Tuple[Any, ...]]:
#     sql, params = "", []
    
#     if kwargs:
#         conditions = " AND ".join(f"{k} = %s" for k in kwargs.keys())
#         sql = f"SELECT * FROM {table} WHERE {conditions}"
#         params = list(kwargs.values())
#     else:
#         sql = f"SELECT * FROM {table}"
    
#     cur = database.cursor()
#     cur.execute(sql, params)
#     return cur.fetchall()
