-- ----------------------------------------------------------------------
-- Target DBMS:           PostgreSQL
-- Project name:          ByteBoard
-- ----------------------------------------------------------------------
--
-- Note: Run this script while connected to the byteboard_db database
-- Command: psql -U kros -d byteboard_db -f database.sql
-- Or: cat database.sql | docker exec -i byte-db psql -U kros -d byteboard_db
-- ----------------------------------------------------------------------

-- Drop tables if they exist
DROP TABLE IF EXISTS comments CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS profiles CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- ----------------------------------------------------------------------
-- Tables
-- ----------------------------------------------------------------------

-- Creating tables
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    hashed_password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL
);

CREATE TABLE profiles (
    user_id INTEGER PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    email VARCHAR(200),
    github_link VARCHAR(75),
    city VARCHAR(50),
    state VARCHAR(50),
    date_registered DATE,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE posts (
    post_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    author VARCHAR(50) NOT NULL,
    date_posted TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE comments (
    comment_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    author VARCHAR(50) NOT NULL,
    date_posted TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_date_posted ON posts(date_posted);
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);

-- ----------------------------------------------------------------------
-- Sample Data
-- ----------------------------------------------------------------------

-- Insert users
INSERT INTO users (username, hashed_password, role) VALUES
('alice_dev', '$2a$10$XYZ123HashForAlice', 'developer'),
('bob_designer', '$2a$10$XYZ123HashForBob', 'designer'),
('carol_admin', '$2a$10$XYZ123HashForCarol', 'admin');

-- Insert profiles for each user
INSERT INTO profiles (user_id, first_name, last_name, email, github_link, city, state, date_registered) VALUES
(1, 'Alice', 'Johnson', 'alice.johnson@example.com', 'github.com/alicedev', 'San Francisco', 'CA', '2024-01-15'),
(2, 'Bob', 'Smith', 'bob.smith@example.com', 'github.com/bobdesigns', 'Seattle', 'WA', '2024-02-20'),
(3, 'Carol', 'Martinez', 'carol.martinez@example.com', 'github.com/caroladmin', 'Austin', 'TX', '2024-03-10');

-- Insert posts
INSERT INTO posts (user_id, title, content, author, date_posted) VALUES
(1, 'Getting Started with PostgreSQL', 'PostgreSQL is a powerful, open source object-relational database system. In this post, I will share my experience migrating from MySQL to PostgreSQL and the benefits I have discovered.', 'alice_dev', '2024-11-01 10:30:00'),
(2, 'Design Principles for Modern Web Apps', 'Creating intuitive user interfaces requires understanding fundamental design principles. Here are my top 5 principles that every designer should know when building modern web applications.', 'bob_designer', '2024-11-05 14:15:00'),
(1, 'Mastering Go Concurrency Patterns', 'Goroutines and channels make Go an excellent choice for concurrent programming. Let me walk you through some practical patterns I use in production code every day.', 'alice_dev', '2024-11-10 09:45:00'),
(3, 'Best Practices for Database Security', 'Securing your database is critical for any application. This guide covers authentication, encryption, and access control strategies to keep your data safe.', 'carol_admin', '2024-11-15 16:20:00'),
(2, 'Responsive Design in 2024', 'With so many device sizes, responsive design is more important than ever. Here are the latest techniques and tools I use to create flexible, mobile-first designs.', 'bob_designer', '2024-11-20 11:00:00');

-- Insert comments for each post
INSERT INTO comments (user_id, post_id, content, author, date_posted) VALUES
(2, 1, 'Great article! I have been considering PostgreSQL for my next project. The migration tips are very helpful.', 'bob_designer', '2024-11-02 08:15:00'),
(3, 2, 'These design principles are spot on. Especially the emphasis on user feedback and accessibility.', 'carol_admin', '2024-11-06 10:30:00'),
(2, 3, 'Concurrency in Go can be tricky at first, but your examples make it much clearer. Thanks for sharing!', 'bob_designer', '2024-11-11 13:20:00'),
(1, 4, 'Database security is often overlooked. This is a comprehensive guide that every developer should read.', 'alice_dev', '2024-11-16 09:45:00'),
(3, 5, 'Mobile-first is definitely the way to go. Your approach to breakpoints is particularly interesting.', 'carol_admin', '2024-11-21 15:10:00');
