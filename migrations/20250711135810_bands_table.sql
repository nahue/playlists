-- +goose Up
-- +goose StatementBegin
CREATE TABLE bands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_bands_user_id ON bands(user_id);
CREATE INDEX idx_bands_name ON bands(name);

-- Create trigger to update updated_at timestamp
CREATE TRIGGER update_bands_updated_at BEFORE UPDATE ON bands
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_bands_updated_at ON bands;
DROP INDEX IF EXISTS idx_bands_name;
DROP INDEX IF EXISTS idx_bands_user_id;
DROP TABLE IF EXISTS bands;
-- +goose StatementEnd
