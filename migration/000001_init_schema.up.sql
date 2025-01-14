CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('client', 'moderator'))
);


CREATE TABLE houses (
    id SERIAL PRIMARY KEY,
    address VARCHAR(255) NOT NULL,
    year_built INTEGER NOT NULL,
    developer VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE flats (
    id SERIAL PRIMARY KEY,
    house_id INTEGER NOT NULL,
    price INTEGER NOT NULL,
    rooms INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('created', 'approved', 'declined', 'on moderation')),
    FOREIGN KEY (house_id) REFERENCES houses(id) ON DELETE CASCADE
);


CREATE TABLE subscriptions (
    email VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    house_id INTEGER NOT NULL,
    subscribed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, house_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (house_id) REFERENCES houses(id) ON DELETE CASCADE
);