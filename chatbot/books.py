from fastapi import FastAPI

app = FastAPI()

BOOKS = [
    {'title': 'TITLE ONE', 'author': 'Author One', 'category': 'science'},
    {'title': 'TITLE TWO', 'author': 'Author TWO', 'category': 'science'},
    {'title': 'TITLE THREE', 'author': 'Author THREE', 'category': 'science'},
    {'title': 'TITLE FOUR', 'author': 'Author FOUR', 'category': 'science'},
    {'title': 'TITLE FIVE', 'author': 'Author FIVE', 'category': 'science'},
    {'title': 'TITLE SIX', 'author': 'Author SIX', 'category': 'science'},
]

@app.get("/books")
async def read_all_books():
    return BOOKS

@app.get("/books/{dynamic_param}")
async def read_all_books(dynamic_param: str):
    return {'dynamic_param': dynamic_param}
