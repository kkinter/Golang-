import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

const Movie = () => {
    // empty obj
    const [movie, setMovies] = useState({});

    // index.js 의 :id 와 일치해야한다.
    let { id } = useParams();

    useEffect(() => {
        let myMovie = {
            id: 1,
            title: "Highlander",
            release_date: "1986-03-07",
            runtime: 116,
            mpaa_rating: "R",
            description: "Some long description"
        }

        setMovies(myMovie)
    }, [id])

    return (
        <div>
            <h2>Movie: {movie.title} </h2>
            <small><em>{movie.release_date}, {movie.runtime}, Rated {movie.mpaa_rating}</em></small>
            <hr />
            <p>{movie.description}</p>
        </div>
    )
}

export default Movie;