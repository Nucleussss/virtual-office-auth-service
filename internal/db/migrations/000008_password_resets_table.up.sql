CREATE TABLE password_resets (
    token VARCHAR(64) PRIMARY KEY
    user_id UUID NOT NULL REFERENCES users(id),
    expired_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)

CREATE INDEX idx_password_resets_user_id ON password_resets(user_id);
CREATE INDEX idx_password_resets_expired_at ON password_resets(expired_at);