CREATE TABLE password_resets (
    token TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    expired_at TIMESTAMPTZ NOT NULL,
    update_time TIMESTAMPTZ DEFAULT NOW()
)

CREATE INDEX idx_password_resets_user_id ON password_resets(user_id);
CREATE INDEX idx_password_resets_expired_at ON password_resets(expired_at);