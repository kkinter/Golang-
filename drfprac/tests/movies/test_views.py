import json

import pytest

from movies.models import Movie


@pytest.mark.django_db
def test_add_movie(client):
    movies = Movie.objects.all()
    assert len(movies) == 0

    resp = client.post(
        "/api/movies/",
        {
            "title": "q",
            "genre": "b",
            "year": "1998",
        },
        content_type="application/json"
    )
    assert resp.status_code == 201
    assert resp.data["title"] == "q"

    movies = Movie.objects.all()
    assert len(movies) == 1

@pytest.mark.django_db
def test_add_movie_invalid_json(client):
    movies = Movie.objects.all()
    assert len(movies) == 0

    resp = client.post(
        "/api/movies/",
        {},
        content_type="application/json"
    )
    assert resp.status_code == 400

    movies = Movie.objects.all()
    assert len(movies) == 0


@pytest.mark.django_db
def test_add_movie_invalid_json_keys(client):
    movies = Movie.objects.all()
    assert len(movies) == 0

    resp = client.post(
        "/api/movies/",
        {
            "title": "a",
            "genre": "b",
        },
        content_type="application/json"
    )
    assert resp.status_code == 400

    movies = Movie.objects.all()
    assert len(movies) == 0

@pytest.mark.django_db
def test_get_single_movie(client):
    movie = Movie.objects.create(title="a", genre="b", year="1998")
    resp = client.get(f"/api/movies/{movie.id}/")
    assert resp.status_code == 200
    assert resp.data["title"] == "a"


def test_get_single_movie_incorrect_id(client):
    resp = client.get(f"/api/movies/foo/")
    assert resp.status_code == 404