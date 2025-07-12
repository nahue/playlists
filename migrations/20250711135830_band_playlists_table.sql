-- +goose Up
-- +goose StatementBegin
CREATE TABLE band_playlists (
    id SERIAL PRIMARY KEY,
    band_id INTEGER NOT NULL REFERENCES bands(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE band_playlist_songs (
    id SERIAL PRIMARY KEY,
    playlist_id INTEGER NOT NULL REFERENCES band_playlists(id) ON DELETE CASCADE,
    artist VARCHAR(255) NOT NULL,
    song VARCHAR(255) NOT NULL,
    notes TEXT,
    position INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_band_playlists_band_id ON band_playlists(band_id);
CREATE INDEX idx_band_playlist_songs_playlist_id ON band_playlist_songs(playlist_id);
CREATE INDEX idx_band_playlist_songs_position ON band_playlist_songs(position);

-- Create triggers to update updated_at timestamp
CREATE TRIGGER update_band_playlists_updated_at BEFORE UPDATE ON band_playlists
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_band_playlist_songs_updated_at BEFORE UPDATE ON band_playlist_songs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_band_playlist_songs_updated_at ON band_playlist_songs;
DROP TRIGGER IF EXISTS update_band_playlists_updated_at ON band_playlists;
DROP INDEX IF EXISTS idx_band_playlist_songs_position;
DROP INDEX IF EXISTS idx_band_playlist_songs_playlist_id;
DROP INDEX IF EXISTS idx_band_playlists_band_id;
DROP TABLE IF EXISTS band_playlist_songs;
DROP TABLE IF EXISTS band_playlists;
-- +goose StatementEnd 