-- +goose Up
-- +goose StatementBegin
CREATE TABLE playlist_entries (
    id SERIAL PRIMARY KEY,
    artist VARCHAR(255) NOT NULL,
    song VARCHAR(255) NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_playlist_entries_user_id ON playlist_entries(user_id);
CREATE INDEX idx_playlist_entries_artist ON playlist_entries(artist);
CREATE INDEX idx_playlist_entries_song ON playlist_entries(song);
CREATE INDEX idx_playlist_entries_user_name ON playlist_entries(user_name);

-- Create trigger to update updated_at timestamp
CREATE TRIGGER update_playlist_entries_updated_at BEFORE UPDATE ON playlist_entries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_playlist_entries_updated_at ON playlist_entries;
DROP INDEX IF EXISTS idx_playlist_entries_user_name;
DROP INDEX IF EXISTS idx_playlist_entries_song;
DROP INDEX IF EXISTS idx_playlist_entries_artist;
DROP INDEX IF EXISTS idx_playlist_entries_user_id;
DROP TABLE IF EXISTS playlist_entries;
-- +goose StatementEnd
