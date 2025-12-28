-- Enable extensions for future proofing
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =========================
-- users table
-- =========================
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    label TEXT,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


-- =========================
-- modules table
-- =========================
CREATE TABLE modules (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =========================
-- topics table
-- =========================
CREATE TABLE topics (
    id BIGSERIAL PRIMARY KEY,
    module_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,

    title TEXT NOT NULL,
    content TEXT NOT NULL,

    is_deleted BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_topics_module
        FOREIGN KEY (module_id)
        REFERENCES modules(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_topics_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
);



-- =========================
-- posts table
-- =========================
CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    topic_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,

    parent_post_id BIGINT NULL,

    content TEXT NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_posts_topic
        FOREIGN KEY (topic_id)
        REFERENCES topics(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_posts_user
        FOREIGN KEY (user_id)
        REFERENCES users(id),

    CONSTRAINT fk_posts_parent
        FOREIGN KEY (parent_post_id)
        REFERENCES posts(id)
);


-- =========================
-- comments table
-- =========================
CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    post_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_comments_post
        FOREIGN KEY (post_id)
        REFERENCES posts(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_comments_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
);

-- =========================
-- Indexes for performance
-- =========================
CREATE INDEX idx_topics_module_id ON topics(module_id);

CREATE INDEX idx_posts_topic_id ON posts(topic_id);
CREATE INDEX idx_posts_parent_id ON posts(parent_post_id);
CREATE INDEX idx_posts_user_id ON posts(user_id);