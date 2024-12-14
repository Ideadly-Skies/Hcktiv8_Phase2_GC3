-- Drop tables if they already exist (for reusability)
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS user_activity_logs;
DROP TABLE IF EXISTS users;

-- Create the users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    age INT NOT NULL CHECK (age > 0)
);

-- Create the posts table
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    image_url VARCHAR(255),
    user_id INT NOT NULL,
    CONSTRAINT fk_user_posts FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Create the comments table
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    author_id INT NOT NULL,
    post_id INT NOT NULL,
    CONSTRAINT fk_user_comments FOREIGN KEY (author_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_post_comments FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
);

-- Create the user_activity_logs table
CREATE TABLE user_activity_logs (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    description VARCHAR(255) NOT NULL,
    CONSTRAINT fk_user_logs FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Insert seeding data
INSERT INTO users (full_name, email, username, password, age) VALUES
('Obie Ananda', 'Obie.Ananda@example.com', 'obie', '11111111', 30),
('Cristiano Ronaldo', 'CR7@example.com', 'cristiano', '11111111', 25);

INSERT INTO posts (content, image_url, user_id) VALUES
('How to be a millionaire', 'http://example.com/image1.jpg', 1),
('Siuuuuuuuuu!', 'http://example.com/image2.jpg', 2);

INSERT INTO comments (content, author_id, post_id) VALUES
('wow obie!', 2, 1),
('🐐', 1, 2);

INSERT INTO user_activity_logs (user_id, description) VALUES
(1, 'Created a post'),
(2, 'Commented on a post');