import pytest

from movies.models import Movie

@pytest.mark.django_db
def test_movie_model():
    movie = Movie(title="어벤져스", genre="action", year="2018")
    movie.save()
    assert movie.title == "어벤져스"
    assert movie.genre == "action"
    assert movie.year == "2018"
    assert movie.created_date
    assert movie.updated_date
    assert str(movie) == movie.title