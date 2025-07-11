-- +goose Up
-- +goose StatementBegin
CREATE TABLE band_members (
    id SERIAL PRIMARY KEY,
    band_id INTEGER NOT NULL REFERENCES bands(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_band_members_band_id ON band_members(band_id);
CREATE INDEX idx_band_members_name ON band_members(name);
CREATE INDEX idx_band_members_role ON band_members(role);

-- Create trigger to update updated_at timestamp
CREATE TRIGGER update_band_members_updated_at BEFORE UPDATE ON band_members
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_band_members_updated_at ON band_members;
DROP INDEX IF EXISTS idx_band_members_role;
DROP INDEX IF EXISTS idx_band_members_name;
DROP INDEX IF EXISTS idx_band_members_band_id;
DROP TABLE IF EXISTS band_members;
-- +goose StatementEnd
