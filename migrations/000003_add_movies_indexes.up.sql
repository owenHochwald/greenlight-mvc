CREATE INDEX IF NOT EXISTS movie_genre_idx ON movies USING GIN (genres);
CREATE INDEX IF NOT EXISTS movie_title_idx ON movies (title);
