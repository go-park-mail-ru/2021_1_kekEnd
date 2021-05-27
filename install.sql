\set ON_ERROR_STOP 1

DROP DATABASE IF EXISTS mdb;
DROP user IF EXISTS mdb;
CREATE DATABASE mdb;
CREATE user mdb WITH PASSWORD 'mdb';

\connect mdb

CREATE SCHEMA mdb;
GRANT usage ON SCHEMA mdb TO mdb;


CREATE TABLE mdb.users
(
    login               VARCHAR(100) PRIMARY KEY,
    password            VARCHAR(256) NOT NULL,
    img_src             text DEFAULT 'http://89.208.198.186:8085/avatars/default.jpeg',

    firstname           VARCHAR(100),
    lastname            VARCHAR(100),
    sex                 INTEGER CONSTRAINT sex_t CHECK (sex = 1 OR sex = 0),
    email               VARCHAR(100) NOT NULL UNIQUE,
    registration_date   timestamp NOT NULL DEFAULT NOW(),

    description         VARCHAR(600),
    movies_watched      INTEGER DEFAULT 0,
    reviews_count       INTEGER DEFAULT 0,
    subscribers_count   INTEGER DEFAULT 0,
    subscriptions_count INTEGER DEFAULT 0,
    friends_count       INTEGER DEFAULT 0,
    user_rating         INTEGER DEFAULT 0
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.users TO mdb;
COMMENT ON TABLE mdb.users IS 'Пользователи';


CREATE TABLE mdb.movie
(
    id              serial PRIMARY KEY,
    title           text not null,
    description     text,
    productionYear  integer,
    country         VARCHAR(100)[],
    slogan          VARCHAR(100),
    director        VARCHAR(50),
    scriptwriter    VARCHAR(50),
    producer        VARCHAR(50),
    operator        VARCHAR(50),
    composer        VARCHAR(50),
    artist          VARCHAR(50),
    montage         VARCHAR(50),
    budget          VARCHAR(50),
    duration        VARCHAR(50),
    poster          text,
    banner          text,
    trailerPreview  text,

    rating          REAL DEFAULT 0.0,
    rating_count    INTEGER DEFAULT 0
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.movie TO mdb;
COMMENT ON TABLE mdb.movie IS 'Фильмы';


CREATE TABLE mdb.genres
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100)
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.genres TO mdb;
COMMENT ON TABLE mdb.genres IS 'Жанры';


CREATE TABLE mdb.movie_genres
(
    movie_id INTEGER REFERENCES mdb.movie (id) ON DELETE CASCADE,
    genre_id INTEGER REFERENCES mdb.genres (id) ON DELETE CASCADE
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.movie_genres TO mdb;
COMMENT ON TABLE mdb.movie_genres IS 'Связочная таблица movie-genres';


CREATE TABLE mdb.actors
(
    id              serial PRIMARY KEY,
    name            text not null,
    biography       text,
    birthdate       VARCHAR(100),
    origin          VARCHAR(100),
    profession      VARCHAR(200),
    avatar          text
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.actors TO mdb;
COMMENT ON TABLE mdb.actors IS 'Актеры';

CREATE TABLE mdb.favorite_actors
(
    user_login VARCHAR(100) REFERENCES mdb.users (login) ON DELETE CASCADE,
    actor_id INTEGER REFERENCES mdb.actors (id) ON DELETE CASCADE
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.favorite_actors TO mdb;
COMMENT ON TABLE mdb.favorite_actors IS 'Любимые актеры';

CREATE TABLE mdb.movie_actors
(
    movie_id INTEGER REFERENCES mdb.movie (id) ON DELETE CASCADE,
    actor_id INTEGER REFERENCES mdb.actors (id) ON DELETE CASCADE
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.movie_actors TO mdb;
COMMENT ON TABLE mdb.movie_actors IS 'Связочная таблица movie-actors';

CREATE TABLE mdb.meta
(
    version          serial PRIMARY KEY,
    movies_count     INTEGER,
    users_count      INTEGER,
    available_genres VARCHAR(100)[]
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.meta TO mdb;
COMMENT ON TABLE mdb.meta IS 'Метаинформация';


CREATE OR REPLACE FUNCTION update_meta() RETURNS TRIGGER AS $$
BEGIN
    IF TG_TABLE_NAME = 'movie' THEN
        UPDATE mdb.meta
        SET movies_count = movies_count + 1;
        RETURN NEW;
    ELSIF TG_TABLE_NAME = 'users' THEN
        UPDATE mdb.meta
        SET users_count = users_count + 1;
        RETURN NEW;
    END IF;
end;
$$ LANGUAGE plpgsql;
CREATE TRIGGER tr_added_movie
    AFTER INSERT ON mdb.movie FOR EACH ROW EXECUTE PROCEDURE update_meta();
CREATE TRIGGER tr_added_user
    AFTER INSERT ON mdb.users FOR EACH ROW EXECUTE PROCEDURE update_meta();


-- TO DO Сделать тригер на пересчет рейтинга в поле rating таблицы mdb.movie
CREATE TABLE mdb.movie_rating
(
    user_login VARCHAR(100) REFERENCES mdb.users (login) ON DELETE CASCADE,
    movie_id INTEGER REFERENCES mdb.movie (id) ON DELETE CASCADE,
    rating INTEGER CONSTRAINT from_one_to_ten_rating CHECK (rating >= 1 AND rating <= 10) NOT NULL,
    creation_date timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_login, movie_id)
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.movie_rating TO mdb;
COMMENT ON TABLE mdb.movie_rating IS 'Рейтинг фильмов';


CREATE OR REPLACE FUNCTION rating_recalc() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE mdb.movie
        SET rating=rating + (NEW.rating - rating) / (rating_count + 1), rating_count = rating_count + 1 WHERE id=NEW.movie_id;
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        UPDATE mdb.movie
        SET rating=(rating * rating_count - OLD.rating + NEW.rating) / rating_count WHERE id=NEW.movie_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE mdb.movie
        SET rating=CASE WHEN rating_count = 1 THEN 0 ELSE (rating * rating_count - OLD.rating) / (rating_count - 1) END,
            rating_count = rating_count - 1 WHERE id=OLD.movie_id;
        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER tr_movie_rating
    AFTER INSERT OR UPDATE OR DELETE ON mdb.movie_rating FOR EACH ROW EXECUTE PROCEDURE rating_recalc();


CREATE TABLE mdb.watched_movies
(
    user_login VARCHAR(100) REFERENCES mdb.users (login) ON DELETE CASCADE,
    movie_id INTEGER REFERENCES mdb.movie (id) ON DELETE CASCADE
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.watched_movies TO mdb;
COMMENT ON TABLE mdb.watched_movies IS 'Просмотренные фильмы';


CREATE TABLE mdb.users_review
(
    id serial PRIMARY KEY,
    user_login VARCHAR(100) REFERENCES mdb.users (login) ON DELETE CASCADE,
    movie_id INTEGER REFERENCES mdb.movie (id) ON DELETE CASCADE,
    review_type INTEGER CONSTRAINT review_type_t CHECK (review_type = -1 OR review_type = 0 OR review_type = 1),
    title VARCHAR(200),
    content text,
    creation_date timestamp NOT NULL DEFAULT NOW(),
    UNIQUE (user_login, movie_id)
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.users_review TO mdb;
COMMENT ON TABLE mdb.users_review IS 'Просмотренные фильмы';


CREATE TABLE mdb.friends
(
    friend_1 VARCHAR(100) REFERENCES mdb.users (login) ON DELETE CASCADE,
    friend_2 VARCHAR(100) REFERENCES mdb.users (login) ON DELETE CASCADE,
    friend_status INTEGER CONSTRAINT friend_status_t CHECK (friend_status = 0 OR friend_status = 1 OR friend_status = 2),
    PRIMARY KEY (friend_1, friend_2)
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.friends TO mdb;
COMMENT ON TABLE mdb.friends IS 'Друзья';




CREATE TABLE mdb.subscriptions
(
    user_1 VARCHAR(100) REFERENCES mdb.users (login) ON DELETE CASCADE,
    user_2 VARCHAR(100) REFERENCES mdb.users (login) ON DELETE CASCADE,
    PRIMARY KEY (user_1, user_2)
);

GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.subscriptions TO mdb;

COMMENT ON TABLE mdb.subscriptions IS 'Подписки';



-- TO DO Сделать тригер на пересчет рейтинга в поле user_rating таблицы mdb.users
CREATE TABLE mdb.users_rating
(
    user_1 VARCHAR(100) REFERENCES mdb.users (login) ON DELETE CASCADE,
    user_2 VARCHAR(100) REFERENCES mdb.users (login) ON DELETE CASCADE,
    user_rating INTEGER CONSTRAINT user_rating_t CHECK (user_rating = -1 OR user_rating = 1),
    PRIMARY KEY (user_1, user_2)
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.users_rating TO mdb;
COMMENT ON TABLE mdb.users_rating IS 'Рейтинг пользователей';


CREATE TABLE mdb.playlists
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    ownerName VARCHAR(100) REFERENCES mdb.users (login) ON DELETE CASCADE,
    isShared BOOLEAN
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.playlists TO mdb;
COMMENT ON TABLE mdb.playlists IS 'Плейлисты';


CREATE TABLE mdb.playlistsWhoCanAdd
(
    username VARCHAR(100) REFERENCES mdb.users (login),
    playlist_id INTEGER REFERENCES mdb.playlists (id)
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.playlistsWhoCanAdd TO mdb;
COMMENT ON TABLE mdb.playlistsWhoCanAdd IS 'Права на изменение плейлистов';


CREATE TABLE mdb.playlistsMovies
(
    playlist_id INTEGER REFERENCES mdb.playlists (id) ON DELETE CASCADE,
    movie_id INTEGER REFERENCES mdb.movie (id) ON DELETE CASCADE,
    addedBy VARCHAR(100) REFERENCES mdb.users (login) ON DELETE CASCADE
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.playlistsMovies TO mdb;
COMMENT ON TABLE mdb.playlistsMovies IS 'Фильмы в плейлистах';


GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA mdb TO mdb;




-- INSERT INTO mdb.movie (title, description, productionYear, country,
--                        slogan, director, scriptwriter, producer,
--                        operator, composer, artist,
--                        montage, budget, duration, poster, banner,
--                        trailerPreview)
-- VALUES ('Чужой', 'Группа космонавтов высаживается на неизвестной планете и знакомится с ксеноморфом. Шедевр Ридли Скотта',
--         1979, '{"Великобритания", "США"}',
--         'slogan', 'Ридли Скотт', '---', '---', '---', '---', '---', '---', '---', '---',
--         'https://avatars.mds.yandex.net/get-kinopoisk-image/1704946/14af6019-b2fe-4e1e-bee5-334d9e472d94/300x450',
--         'https://kinopoisk-ru.clstorage.net/E3Y0U4315/f882714f3/JQb5hfuBU5QkfR6H5sDiqWH4CdAOpqGh5AUpIFJuMUGWopdsc4iigraQMDvjNkTICX7J3xVeK_VrWrgjPNqrlmALqJy2hRKhTQhu5t-ku_QU2c0UnnMVQUycuYCmh1BUhv_iLRaJqrvTnVxC6ifhcEGWDVqAjgfY3g9BZlSSplf1Sl3Z190y4Mh9y-faZGYlbfGJsHrf3RU8MFEw5h8klIt7lkES6bq7e0FppppDjcEJSvHiwCpzxHJHtTIsVq7fxcZxZSddQsQIBfsr0ozC_VXNYWxK3xklyBAVWJ5jzNxH9-aFovVSj_OZ3bayy6ERobZp1oAaM4TaDwk61PJSK0nixcVXaaq0BEn7H2bYBgl9kQHsDk8hEfTQiXhmrwhoS1-eyXYBjpPbbZVqgsNdNcV-HXIMdoNUag_lhghqhn-tNjFx0-HG_DxB369K9JpxKfEJdPK3GXEQAEV05heU7Lvv6jV6UYq_S-GZ5nJXNSl1upnW2N5HMJIzbdZMgmIvPb6twSdVPmQkIfOzGugKkSk9DeQiwzX5IFx9DBaHiBCbLwbhsqE-A9tJ3Zoeo5GJnd6Z9piSI8wSrxXOIPLS6wE2bc2nxd7MXN3n_-JManUpLWHUSs-FiaTEtaiyNwSca5OG0YrppiP7bTUq_qcFXdlWEfJUTovY7kPJstSyjsNtwjEZv21GhJQtn98mlFZl6cklrFZn_ZFstCEUhg8U3BfLjsmKldIj1-2tdqrLJdUlojHWcCK34NIPee5UzuaTKT7FWTuZxhAg3VfnulTC6UkFnai2q1G1WASh_LYzJIB_S2Y9GuVOF__hzWYmoxndUSLBGtiaE6SisxFaXGqWi_3Kyc3v4d7EwLHXNxLwmn2xCZkw2i8RpdjwhTS2uxSQ-4veYbYhVk_bEd3echshHSky-SpEessoHovVBnBe2hcB2qkBV-0u-BxFF-_i1EIxfQmh5H6TidlE9MWkplecKB___sk69VZnMxkxQkonGcF5ss3atLbrRPLT0drwkgYvrTq1ed-ZGow0YUcbvhTejVHRcSxK_9m5EEjZNOb3yLjHX3K59lFWJ9-1vWoe4w0ZgeKF4vS261iOv4VKiOqKhzWCZYXrmVaMzEUHe0Z4GvmlTfmAMmeBeaTw3TA61xhAX_9S6RLV9heT4a2aNltZQfFmTfp4pptI2qPtfiTK6p81GjWxd4XeOIBRyxdiyA4teT0hpEbLTV1MQDF4tlsgaM-jxg2qqb4zY-3NVqp7tbkJFpEWwDqPbBp_iXJgfm4HqR41ydcdtjQ4BUOjBuBq7cGB5XjS_40htMQZ1B5rTLQXH_qd5mHiHz9NPRIeZ1EFfTYJ1vh2h0xe08lueA56f3F2tR13AaK8vPUPN5boVrExQQFIhqsJdei81Yw2k0wQM4N-SSoFPm9Xlbne-juVRVn-cQ7QjifU-iu17njqNp-lkuXh36EKvIzlQ3PC0JpdwQGRWH5vCXkYEDX0UqtoDGtPhjW6idbbk-k9xq7DTTnRpkWuEE5_xBID7SZ8TrIfJXaVTceFqmi8_UM31sCGtU2xEcAK59mFRHjJ2JofLERvp5bhsi0mo_v5TaaK7wk1ESadkoTCh1SOl0lOCJK2bzm6xV2ncRpoHCELJ7Z0FhW1zQFkun_dMdhsqZQ6a5Akfy-ubT5p9stf1dVCSk9BHRkSDWZwAmfA8jcNPnAy8mMR0vkxP_kWqFjlE_9CbFIJ0fGpSMLLPSGUhNH0rjsstPsXHi2avYLDq529Ig4rQcEVto3ucAY_7IrbUUo03qYT6XapGUsdGhBULf8jkhyysfEtPbSe40HVfMQJ-FKznDAT61qhFpnmvx8x1ZrypwWF5eL90vBan1wqm1mGqDY-RwHWaUmrjX50rIkfXwr0noFRQSFciq99lbiYRdwmN7BU19cmKf4F2iuvhXWegsOlAXkWQV5gJkv0ZlMBovTCvkeZYkkZaw0KJDxRV3tSzMpdxd3hJCJPSSFsMIlIkpOkrFfjXrXGBaZLyzXJRmbDJQXFNknaDK7zwHYrtX44XnZP-Zo1nZ8dDqTEoVuHBmRCpS2pHXQ218m1-PjNVK6DLMDPU3btoo32AyvZGWr-G421JbJJGhS-h1galz0CvIJqHw3CPT3f9UoUrPHPFz70Tn1pxTUMQnO5MeDwtdSmVyw0wwtezTYlvscnlbFmlps1KclabYpQSjOgPt-RPozuGs-tPh1xJ4n6dDhBe6uWeMrtxYlh4DZLKXWQ-HWQRicENH_PSl163bJTn8V5xkJXkWUd5kFKfALHaLILASrAzr5TPeapFb9tCjwUWRc7opDiLRUh5Xye3919oIR1vCbrbMgDw24hvgm6CxONlfLKby1BXU7FQsCKl4SWx-VqjEb6W2XGcWFzYRZINM2LA6aECl3dvTlE3ru1dbQw3VReN6gkV4NiQTJt5uf_5enmtj-9XYE2UfpEdgM8XgudXlxuAocFPkX5lwXaAFRlY6_C0IKxnalxxEqfLQHU6D1IRjfkBB8LGqVugcofq5lNarY3AYWF6hmOhDoTUO7TPe4k9qZjpQYZdSfdQjyEfU_nZtx2HS0lKaCWpyGl2IxdiEZ3FKSXy8KFSi2iI2s1YeZq2zFNqeppDrh6h0iGP5UqzDbSG0UKXdmjkXooUGFbo-ZI4h1xuXEkOmsJZUycCawqayhcX5umHUY98k_7Bb1mgr85HRV-MZbEVstgdg8VhrjGoiPtjk0Vl_3aIMBBh3_GbIoFabFVsCpPmY2UBBXUFhewXGeXLqlO0aIne-1NXma3zcmZTjlSSBq3bB6HWXr0FhpvmbalUXMZRuRU3XcTEpy--R1dVbCWy_UZVPilgLLvJOzva0YVSjGqH18dKcYKw5m9Bf6x-pTWJ-ieu51WIAJ-X7naeQ27RbqgLOEDL-L81pExhWlIKk9JuWCciQiG59Ckc0uSNfKVJgMnefkiAt_9nUl6nVY8vvM0NnsA',
--         'https://kinopoisk-ru.clstorage.net/E3Y0U4315/f882714f3/JQb5hfuBU5QkfR6H5sDiqWH4CdAOpqGh5AUpMCYHMTyD-989NpHas7rlbX6SQwnFJE-UnxlPftgzVp1zDZvnkzQK7JHuiFPIBGAbpt-Zg-gkxTiJQ6LRsM0RXFxjZlz1lzsyjao9Fjs_de3eGi9B0Hle_fZgXhPAdjLlKjwuR8vFimkJb-EiSADxV4PC2EZtYTk9wA7b0WXgzEVEooOQELN_lmkG9UJLk51dnoJbgd39KpHyAE4TaKoHAVpYQvaDGXbRScuVPuRIBQP7BtDOXe0VgTCOq_WlmPyJTEonSCiPl0pZBjFyv7PZNeb-s8kRcVpVQhymJ-BysxFy0G4mK-Xu5Qn71X5MALkLA6IsOtHJoSFAukOFqdT8vfge47wA0_fKHZYJRoNzRRUeMnd9HdkiGZboejvIhvvNQtiSAk8tmkGFl2Wu_ABVt986-G41HVWVxDIvzZHc2A3UMn_EVLP_UuFuMSZno4kZnorr-e1pVuWauH4XWAr3QbJM-h7_6R7lUT-hygQ83eeb0hw-pZXJJTCG11WJvPQRhC7vJFR_-8IxGh2uV_-xYW4a9wXNqboxUkjKq8iisz3euF6-F40exRFzxTaEoKX_96Ks-pkdCSmwFjf98UDgXXyKj8BUb5fyWRat3tNnecXKqnuJPRU6AWoAUovohluNPrzKalMFlsHdr20ieEx5gwPi8NL16VX9qL6vtTmwmH24UrPQlIvTikW-1cYbF-15_pJryUFNMhlqfCaLxAbD0WrQ6uKv8bbl-cNRGkQAyd-DnpiCsRWhvSxKLyGNQFBFJJInXDRHa46lcnniL6dtkc6uW5Upzdrt-gy6v-wKo8HmuNbq23FGKVF79V40vKFrizrommXhrSn4Mjf1bSzQlYw2f8jMS28WyfY58q9TSVnOJmuFrQ1isVbIoufI-rN5sgDuKqNhfhnNmy3SiIRlN6cOpAaZ8c3lQD7HybHYEE18EqeEAEtXwm1KoY4zVwnJ3srjPUl5QhnaHKLPIPJf5Yo81vbz4UrpPVcNvmTcYesnwng-NRHRnchK872Z_EC5INI7OCyThwpZJvHuZ-sVWZ5qt62R2c5pFriij8xe083e-MIuC7EC0X1XDaIYsDV7X7r0lq2pAWH8Sr-9YdgA2di-_0zYDw-mIb6pLtNTEV1CSmdVCXnuOfI8Ar-ACsM99kCWdns1ysnxR32yTKxdT_OalI6tMVFVYFY3CS3MzLX8DuuYBH_XglUSZQo74_0VzsZffZklet1KQEqbcAaj8Wpgeo6DRRYlSdtploxwOUO3LhAWMTVRLcDOXwWVmEQBmCaPWLzDE17BJqV2w2fVuWb2M6FBmUZNBogWtyymU7XefJ4y92WO5XGXYbbI3Hlfr14Ebuld0flg0kuNEWgIlQguswRMA_dulXIhIp8fGeFODjMFZQXCmcrsysdEftd5OiBactOt9j1Zb8EubCQF36-6SI49uYEFyHLjjSF4RNFcFn_ovENnfm22IS5vs_mZKjYXGT3JOuVaYCJzgAJTYW7Ygg5b9cKdma-ZPoQMXRerHswOvV3xqdBWQ1kRYESVSAZjADDz5-YZPvHSM9sFteKCU1E5ISoxUsTSC-gSIwFK9MYCm3UaoQ0jYa4YmPl_38LIfqGRobmwovNZsbwMhSiy86DIj_dCqab1Zq_PZflC9u8xKakSvd6AAmNMPrvlilSOKpNBilX544E6ZDi9D6dijHKJ-Z3VKCr_mfV4FF3cqre8rLNfbtESFXbjJx2Z1qZToa2Rov16VHZruHbThc4wjvaf5Qrd-efZFhzU4XvjjtgCcV3N_VzO8yH5sPiBDNpXBIxvy5KNOmmCC2fFlSou4yVFbeZx9nASFwzauz0yvMqyb7F64Xm7eaa8lOm3f2ZAVpn9Da28XpdFARQY_ZQyezQsA9d6mXZVws87ibFeqs9BgVGa-R7sLoO8bhs5QthqNvNFxm3px60O8FyxkyOSwFYBSS39fN7jFZHMUNnMCi_ouJ8XAjGWYXYbk0Ul6g7buQFl4mUm7FLj2N6n4abY6jJPZc7phU8VOuAkBU_vDgheYbFReYjO55VpPFwlmKKnEFDr61IlDuHij1sBOdYeU9WZ1co9QmQCqzgyd80-AEKCr-HOKZ1fYaKMmI0za9IUDpXpWdnIJqMlAWzItaAyq8gUh8MqUaqRZpdTebneylMhlY3iHdbMSm80ft_BVoD6HkMJ6rnZq9VaqNAhD1u-ZN41FXmVMFoTRZXcfAkIvi9YuMuXxiWSASLnW7n9Prp7ISlJ9o2aNEb7jC4XYYJMXlKXtcZ59eMhkiQEsRsXnsBCpc3N8ai-4w25xBCZPFYHmGhjE1qNBvUq1ye50V52E91VRdLxXuBOowBm-1UKdOJ21x1CcUlrcX4AyFVbWxaESv3tFYVksv95mVCMoThC7-ig_89izWKdIsOTETkmqtcxAQXekdKEEk_sDodBdiRyagtl1snNl-XGyAQtb4s-fJadFSEdgNYzMfn4ZA1cFmcE4OuH4llGBVajS_ElPqqbEUmNpnWOaD63uHIjzXYszrIPuZ69Ddv1qnjcjd_zpthyPS19kTAOqw0p4EhF-BqTqFBn34aFfgnyry-R5T7qa7HBTX5VqsRWi3jeD0GqwP56I7nuPTGbYbIQMCUbG2asCt0hOT20QpMZ_fxcAXiOB6gM-4cCKbIhMjs_xcFS9ldJCR0azabUBufo7tPBQqT2Kp8ttqVNty2a4AClt2-W3DJ1pSnxgC4zEW3cgN1Yqm-wFPOjljmWsdrjp9m5borPSTERknmuOFaPaAYj-aasAv4THb5hwftRloiI6UsjRmR-AZ3BtWTKr9X5QHCxjFpbTGAfo5aFEt1OI1tp7cpyW_m57frFqthet0z2R2HK2FaKj602yR03wRIItC1n91IATiHxHemsllORgXwEjXw6MyRMx59uOZZh7hc_RWX-eq-xJc0u5RJ80qs0kpeE'
--         ),
--        ('Назад в будущее', 'Безумный ученый и 17-летний оболтус тестируют машину времени и наводят шороху в 1950-х. Классика кинофантастики',
--         1985, '{"США"}',
--         'Семнадцатилетний Марти МакФлай пришел вчера домой пораньше. На 30 лет раньше', 'Роберт Земекис', '---', '---', '---', '---', '---', '---', '---', '---',
--         'https://avatars.mds.yandex.net/get-kinopoisk-image/1599028/73cf2ed0-fd52-47a2-9e26-74104360786a/300x450',
--         'https://kinopoisk-ru.clstorage.net/E3Y0U4315/f882714f3/JQb5hfuBU5QkfR6H5sDiqWH4CdAOpqGh5AUpIFJuMUGSio9Ic4iigraQMDvjNkTICBLdyzlDZqAyCoF3DMPzgzgO9dy3wQv1XThu5t-ku_QU2c0UnnMVQUycuYCmh1BUhv_iLRaJqrvTnVxC6ifhcEGWDVqAjgfY3g9BZlSSplf1Sl3Z190y4Mh9y-faZGYlbfGJsHrf3RU8MFEw5h8klIt7lkES6bq7e0FppppDjcEJSvHiwCpzxHJHtTIsVq7fxcZxZSddQsQIBfsr0ozC_VXNYWxK3xklyBAVWJ5jzNxH9-aFovVSj_OZ3bayy6ERobZp1oAaM4TaDwk61PJSK0nixcVXaaq0BEn7H2bYBgl9kQHsDk8hEfTQiXhmrwhoS1-eyXYBjpPbbZVqgsNdNcV-HXIMdoNUag_lhghqhn-tNjFx0-HG_DxB369K9JpxKfEJdPK3GXEQAEV05heU7Lvv6jV6UYq_S-GZ5nJXNSl1upnW2N5HMJIzbdZMgmIvPb6twSdVPmQkIfOzGugKkSk9DeQiwzX5IFx9DBaHiBCbLwbhsqE-A9tJ3Zoeo5GJnd6Z9piSI8wSrxXOIPLS6wE2bc2nxd7MXN3n_-JManUpLWHUSs-FiaTEtaiyNwSca5OG0YrppiP7bTUq_qcFXdlWEfJUTovY7kPJstSyjsNtwjEZv21GhJQtn98mlFZl6cklrFZn_ZFstCEUhg8U3BfLjsmKldIj1-2tdqrLJdUlojHWcCK34NIPee5UzuaTKT7FWTuZxhAg3VfnulTC6UkFnai2q1G1WASh_LYzJIB_S2Y9GuVOF__hzWYmoxndUSLBGtiaE6SisxFaXGqWi_3Kyc3v4d7EwLHXNxLwmn2xCZkw2i8RpdjwhTS2uxSQ-4veYbYhVk_bEd3echshHSky-SpEessoHovVBnBe2hcB2qkBV-0u-BxFF-_i1EIxfQmh5H6TidlE9MWkplecKB___sk69VZnMxkxQkonGcF5ss3atLbrRPLT0drwkgYvrTq1ed-ZGow0YUcbvhTejVHRcSxK_9m5EEjZNOb3yLjHX3K59lFWJ9-1vWoe4w0ZgeKF4vS261iOv4VKiOqKhzWCZYXrmVaMzEUHe0Z4GvmlTfmAMmeBeaTw3TA61xhAX_9S6RLV9heT4a2aNltZQfFmTfp4pptI2qPtfiTK6p81GjWxd4XeOIBRyxdiyA4teT0hpEbLTV1MQDF4tlsgaM-jxg2qqb4zY-3NVqp7tbkJFpEWwDqPbBp_iXJgfm4HqR41ydcdtjQ4BUOjBuBq7cGB5XjS_40htMQZ1B5rTLQXH_qd5mHiHz9NPRIeZ1EFfTYJ1vh2h0xe08lueA56f3F2tR13AaK8vPUPN5boVrExQQFIhqsJdei81Yw2k0wQM4N-SSoFPm9Xlbne-juVRVn-cQ7QjifU-iu17njqNp-lkuXh36EKvIzlQ3PC0JpdwQGRWH5vCXkYEDX0UqtoDGtPhjW6idbbk-k9xq7DTTnRpkWuEE5_xBID7SZ8TrIfJXaVTceFqmi8_UM31sCGtU2xEcAK59mFRHjJ2JofLERvp5bhsi0mo_v5TaaK7wk1ESadkoTCh1SOl0lOCJK2bzm6xV2ncRpoHCELJ7Z0FhW1zQFkun_dMdhsqZQ6a5Akfy-ubT5p9stf1dVCSk9BHRkSDWZwAmfA8jcNPnAy8mMR0vkxP_kWqFjlE_9CbFIJ0fGpSMLLPSGUhNH0rjsstPsXHi2avYLDq529Ig4rQcEVto3ucAY_7IrbUUo03qYT6XapGUsdGhBULf8jkhyysfEtPbSe40HVfMQJ-FKznDAT61qhFpnmvx8x1ZrypwWF5eL90vBan1wqm1mGqDY-RwHWaUmrjX50rIkfXwr0noFRQSFciq99lbiYRdwmN7BU19cmKf4F2iuvhXWegsOlAXkWQV5gJkv0ZlMBovTCvkeZYkkZaw0KJDxRV3tSzMpdxd3hJCJPSSFsMIlIkpOkrFfjXrXGBaZLyzXJRmbDJQXFNknaDK7zwHYrtX44XnZP-Zo1nZ8dDqTEoVuHBmRCpS2pHXQ218m1-PjNVK6DLMDPU3btoo32AyvZGWr-G421JbJJGhS-h1galz0CvIJqHw3CPT3f9UoUrPHPFz70Tn1pxTUMQnO5MeDwtdSmVyw0wwtezTYlvscnlbFmlps1KclabYpQSjOgPt-RPozuGs-tPh1xJ4n6dDhBe6uWeMrtxYlh4DZLKXWQ-HWQRicENH_PSl163bJTn8V5xkJXkWUd5kFKfALHaLILASrAzr5TPeapFb9tCjwUWRc7opDiLRUh5Xye3919oIR1vCbrbMgDw24hvgm6CxONlfLKby1BXU7FQsCKl4SWx-VqjEb6W2XGcWFzYRZINM2LA6aECl3dvTlE3ru1dbQw3VReN6gkV4NiQTJt5uf_5enmtj-9XYE2UfpEdgM8XgudXlxuAocFPkX5lwXaAFRlY6_C0IKxnalxxEqfLQHU6D1IRjfkBB8LGqVugcofq5lNarY3AYWF6hmOhDoTUO7TPe4k9qZjpQYZdSfdQjyEfU_nZtx2HS0lKaCWpyGl2IxdiEZ3FKSXy8KFSi2iI2s1YeZq2zFNqeppDrh6h0iGP5UqzDbSG0UKXdmjkXooUGFbo-ZI4h1xuXEkOmsJZUycCawqayhcX5umHUY98k_7Bb1mgr85HRV-MZbEVstgdg8VhrjGoiPtjk0Vl_3aIMBBh3_GbIoFabFVsCpPmY2UBBXUFhewXGeXLqlO0aIne-1NXma3zcmZTjlSSBq3bB6HWXr0FhpvmbalUXMZRuRU3XcTEpy--R1dVbCWy_UZVPilgLLvJOzva0YVSjGqH18dKcYKw5m9Bf6x-pTWJ-ieu51WIAJ-X7naeQ27RbqgLOEDL-L81pExhWlIKk9JuWCciQiG59Ckc0uSNfKVJgMnefkiAt_9nUl6nVY8vvM0NnsA',
--         'https://kinopoisk-ru.clstorage.net/E3Y0U4315/f882714f3/JQb5hfuBU5QkfR6H5sDiqWH4CdAOpqGh5AUpMCYHMTyD-989NpHas7rlbX6SQwnFJE-ckwl_brUGDoguZZf20ygPpcCz3E_xSHQDu4-NiqFgyHX8Cv7U2KVpWF1CXlVYJosy9aoxYntLGdXmBstB0RxK9epkIm_YHlPwWiCGRjYhuilBI1k2kPR9z7u-hMYloY2R5DLjsfEoRIVEWp-wlMOvdlFOnaovO7U9VkZTNRERzoGGYEJ_2LaPxb5Q4iqHaWbV-WP9QoxYNTvvxkDOrZEBvVjCY8HV6Dy1iFJ3FEz7k56Nfp1uH8-VeT4-L91Z3ULxQtBel-w-V3GueGoGV8GaTc0jzQLM8H2H5z7kMlkdJQn4slcppeRwtbzmI9C408_-DToNVivzVeUexuMZ7dHqiQ4EqkvwFqM5ckhi-nOlUjlpr6GyHEB9a1vifOYN-fH9TDbfRe3ceJEMyg9MwIev9pXG9W5LF4UpEkZbhWkhWv3yCPpP3IYvNf649pJvFZa9zXsJdni4QeMLppQCXWl5YfzCa711xBi9EJoT3CCHY_IFFoFCwyfZEWq2y5mVAZoRJsAK-2AWh3GC1AI2z_3yve07RRKEON2bE8rksplV8aHwQvtd3bzkqVxit7zEh3OeNX6N8rOjQdnOEnsVGfEmkRb4QmNANqOZMjQGohu5ejXp95m6kMQxR28-pO6xOQX9JFpTxZV0FNF8pm-A1EeX2k1iJYqrazFNciZDBVmNfpkO-D4XQBojAW5gaoKTRY4VzdP1hqj4ffczvtiG4X35CWTep0UBwOQZRDqvFFjnW2JJgukmj1-BzZoWfzUF5f5x-mhOi3QyL2F-7AK-mzEO5QF7TSLsiMGfh7Z89vmpDQXwCt9d1SCImZSSC0zMH1dm0e5tZp_fdelSFvcFFWE-yabEipMsFt9xxri6hltJHt0x5636YDT5W9uaSLplVR1lPLLTren8fFlMYi-UgNNXXgVK0f7jQ3GpwgYbja2FSukOSF6TBP7XnVqAhr6HGZ7pwRdh2gzYoV8HGoRmXfn9eUQ6p5md1FgJuD7vCDz_j47Nfr2ugxfNtVJGu9k9XeplfoT6k0QSexFy1EKqX-HOoflXYdoQpM0Ll2L86vVhRam4DqfVnSx8SdjGg8xICxMGYQYl9kOjdbFWmpsJxcVKRS5gfjN0Xi8Bgvz6_geRSmnh23GqAPDRY6PO3IrtYd35jJK7XSlgaIW04jPYnNdj3kVyiTpnS8VdHhYXMe1VFtHK2AJ7UK4jYU5g2hL_aTq1DWPtviQwDQevimgOdf3Z-fQyIzUl2DwNAIYbvFxv3xqZ5r36G7NBdbK-J10xjartWpTKJ3zyg5EK1Mb2Qx0aLc1bobYEdKFHs5IYGg0lsXkgkj8hrVzMQZQWE4AAnx_-qbLpfk_vObnqlt9dlak2aY5YrvsMmlsVxjCaMgM50lUVc1kWnNBZOzOS_Fbt8VUp3Dqfia1s3A3QQitM7G9fbrlKLX5DH5VZkvLneYnx-pHyyCITuF4nkd5kYup_sYphtbOZTow4cWP7lljSbXGxWXAiuyl5XMQNlFY7UATj7-4hPqWuv0P9pb46Uz3B9RKBJsCG48A2N-G-QE6uc3EKuYknFbYcpOXHk-KE1h1tfQlgQk-ZefwYRYQ2j8CkG5P-hY49qgvf6cXymieBoeWauapMwjOokht5WoDu5lt5Pil909VWiNhFg-OaJJIRRRU1DNrHlbm43F1cwpeEuH-vVqn2iUobkwG9kg53PTFhognq6BZHoGZTETrEiuaHdZqp9dPRDqSgqd-X3sjGYb2xZSSuI5kBtBSxgBLnZABfc8JVqqE273tBZZ7y_421iV5NZmQyI9zS_3mCOAaiw4XO2clTja4UAOnXW0IgXjVVEaV0TrP9ZUywUfyKD0gw_x_evb7tCq-_HSm6hnuh0U1iMe6Mrh9IYkvZhkhiAkcZOmVFw_F6vEwhj38e1N41zaWFJI4ziTXcaBnY0jcc7GuDHsUWDT4ba7XlLjLftSnNVklytK5jKAb7ZV6sYoJDpRptwa95wohcWTuj0kgWPa1d-aB6I421JJgVJIaflBSD9-KVApW-j_99oTIOzz1FVeZhKtAmM2DmF7VyNLoq80WebQG3abYQMOWz31aUCm1ZBfEAOsvJBUzIgbS-D5jMx5vK7XYxzgvnddmyBhs9sVm-SQpEjnuk6lsdflw6km-pdkmR850C6BStH-Nm-Hq9-fnRTMK3eWXYeDUIFoMcXGvXngECCV5Pl30Z9uZrFbHlel2aCHZ3MFIL1d6I9jYjfcplUd_V9iCYeY_3KtjeIWkhZShaU4kt9GBZmCJrNJy7fxqdqp2qR6cBGdqGp31NmXZ55syif2jeQznqAM6KBz1i4VljXabMvLVrt2ZQmikxAb1cll-VWdT0xaAmf9zsc-PGper5wk-ztbEy_nu5oc02dYZAxiOEMitF_nyeGhvhGnXh56EydHR5E4O2eGL1UfmJxHI7WRG0XC0MQitUADP3jiV-3Vo7021RLuZ79YGFvg1iHCoPfGZX4XJ8lqbD5cY9lSftIhjEobMzzuDGEfHB1UjC48EtZEQBROYnoKyDe9ZBouVWn98JMe7mOwUhDX7VQjiGZ0Cm-83-oHqWC8nGTRUbrbYArE0b9yYgsmkRzZHkRq_5ObBYFQBmszSs3-eOxQ4pfl9LGWXKiic52cUusdo0ljcsNssRfkgenlt1UhWNZ4H6KFx9m1tS0MJRuUmBKHLDWTEgeMncRpdctMfvqlEeDe63k4F5srZbodn9IjluPHpnRLYj4UasFmqP-WIdSevNhiQ09denHgB6Hc1xaWyWJ8X1tOQ5sJJnaEizA6pRoomCI1N9yeYSozVpdd5R0jiab3yS04XewGI--2XSleE3ARagtMkTi8oUHi3tHbUwXns5sczYTYxiBwAgn9uWqR4NPoNnGeVuJqvBIen-hfKAPuNg6rdU'
--        );


-- INSERT INTO mdb.genres (name)
-- VALUES ('comedy'),
--        ('drama');


-- INSERT INTO mdb.movie_genres (movie_id, genre_id)
-- VALUES (1, 1),
--        (1, 2),
--        (2, 1);


-- INSERT INTO mdb.actors (name)
-- VALUES ('act1'),
--        ('act2'),
--        ('act3');


-- INSERT INTO mdb.movie_actors (movie_id, actor_id)
-- VALUES (1, 1),
--        (1, 2),
--        (2, 1);


-- INSERT INTO mdb.users (login, password, email)
-- VALUES ('user1', 'user1', 'qwe1@mail.ru'),
--        ('user2', 'user2', 'qwe2@mail.ru');


-- INSERT INTO mdb.playlists (name, ownerName, isShared)
-- VALUES ('Плейлист 1 usera1', 'user1', True),
--        ('Плейлист 2 usera1', 'user1', False),
--        ('Плейлист 1 usera2', 'user2', True);


-- INSERT INTO mdb.playlistsWhoCanAdd (username, playlist_id)
-- VALUES ('user1', 1),
--        ('user1', 2),
--        ('user1', 3),
--        ('user2', 1);


-- INSERT INTO mdb.playlistsMovies (playlist_id, addedBy, movie_id)
-- VALUES (1, 'user1', 1),
--        (1, 'user2', 2);



--        ('Титаник', 'Гибель легендарного лайнера на фоне запретной любви. Великий фильм-катастрофа — в отреставрированной версии',
--         1997, '{"США", "Мексика"}',
--         '{"Мелодрама", "История"}', 'Ничто на Земле не сможет разлучить их', 'Джеймс Кэмерон', '---', '---', '---', '---', '---', '---', '---', '---',
--         '{"Леонардо ДиКаприо", "Кейт Уинслет"}', 'https://avatars.mds.yandex.net/get-kinopoisk-image/1773646/96d93e3a-fdbf-4b6f-b02d-2fc9c2648a18/300x450',
--         'https://kinopoisk-ru.clstorage.net/E3Y0U4315/f882714f3/JQb5hfuBU5QkfR6H5sDiqWH4CdAOpqGh5AUpIFJuMUGaho9Ee4iigraQMDvjNkTICBO4hxleIoVaGoF2YbP3qnAfqcSumQa9dHxu5t-ku_QU2c0UnnMVQUycuYCmh1BUhv_iLRaJqrvTnVxC6ifhcEGWDVqAjgfY3g9BZlSSplf1Sl3Z190y4Mh9y-faZGYlbfGJsHrf3RU8MFEw5h8klIt7lkES6bq7e0FppppDjcEJSvHiwCpzxHJHtTIsVq7fxcZxZSddQsQIBfsr0ozC_VXNYWxK3xklyBAVWJ5jzNxH9-aFovVSj_OZ3bayy6ERobZp1oAaM4TaDwk61PJSK0nixcVXaaq0BEn7H2bYBgl9kQHsDk8hEfTQiXhmrwhoS1-eyXYBjpPbbZVqgsNdNcV-HXIMdoNUag_lhghqhn-tNjFx0-HG_DxB369K9JpxKfEJdPK3GXEQAEV05heU7Lvv6jV6UYq_S-GZ5nJXNSl1upnW2N5HMJIzbdZMgmIvPb6twSdVPmQkIfOzGugKkSk9DeQiwzX5IFx9DBaHiBCbLwbhsqE-A9tJ3Zoeo5GJnd6Z9piSI8wSrxXOIPLS6wE2bc2nxd7MXN3n_-JManUpLWHUSs-FiaTEtaiyNwSca5OG0YrppiP7bTUq_qcFXdlWEfJUTovY7kPJstSyjsNtwjEZv21GhJQtn98mlFZl6cklrFZn_ZFstCEUhg8U3BfLjsmKldIj1-2tdqrLJdUlojHWcCK34NIPee5UzuaTKT7FWTuZxhAg3VfnulTC6UkFnai2q1G1WASh_LYzJIB_S2Y9GuVOF__hzWYmoxndUSLBGtiaE6SisxFaXGqWi_3Kyc3v4d7EwLHXNxLwmn2xCZkw2i8RpdjwhTS2uxSQ-4veYbYhVk_bEd3echshHSky-SpEessoHovVBnBe2hcB2qkBV-0u-BxFF-_i1EIxfQmh5H6TidlE9MWkplecKB___sk69VZnMxkxQkonGcF5ss3atLbrRPLT0drwkgYvrTq1ed-ZGow0YUcbvhTejVHRcSxK_9m5EEjZNOb3yLjHX3K59lFWJ9-1vWoe4w0ZgeKF4vS261iOv4VKiOqKhzWCZYXrmVaMzEUHe0Z4GvmlTfmAMmeBeaTw3TA61xhAX_9S6RLV9heT4a2aNltZQfFmTfp4pptI2qPtfiTK6p81GjWxd4XeOIBRyxdiyA4teT0hpEbLTV1MQDF4tlsgaM-jxg2qqb4zY-3NVqp7tbkJFpEWwDqPbBp_iXJgfm4HqR41ydcdtjQ4BUOjBuBq7cGB5XjS_40htMQZ1B5rTLQXH_qd5mHiHz9NPRIeZ1EFfTYJ1vh2h0xe08lueA56f3F2tR13AaK8vPUPN5boVrExQQFIhqsJdei81Yw2k0wQM4N-SSoFPm9Xlbne-juVRVn-cQ7QjifU-iu17njqNp-lkuXh36EKvIzlQ3PC0JpdwQGRWH5vCXkYEDX0UqtoDGtPhjW6idbbk-k9xq7DTTnRpkWuEE5_xBID7SZ8TrIfJXaVTceFqmi8_UM31sCGtU2xEcAK59mFRHjJ2JofLERvp5bhsi0mo_v5TaaK7wk1ESadkoTCh1SOl0lOCJK2bzm6xV2ncRpoHCELJ7Z0FhW1zQFkun_dMdhsqZQ6a5Akfy-ubT5p9stf1dVCSk9BHRkSDWZwAmfA8jcNPnAy8mMR0vkxP_kWqFjlE_9CbFIJ0fGpSMLLPSGUhNH0rjsstPsXHi2avYLDq529Ig4rQcEVto3ucAY_7IrbUUo03qYT6XapGUsdGhBULf8jkhyysfEtPbSe40HVfMQJ-FKznDAT61qhFpnmvx8x1ZrypwWF5eL90vBan1wqm1mGqDY-RwHWaUmrjX50rIkfXwr0noFRQSFciq99lbiYRdwmN7BU19cmKf4F2iuvhXWegsOlAXkWQV5gJkv0ZlMBovTCvkeZYkkZaw0KJDxRV3tSzMpdxd3hJCJPSSFsMIlIkpOkrFfjXrXGBaZLyzXJRmbDJQXFNknaDK7zwHYrtX44XnZP-Zo1nZ8dDqTEoVuHBmRCpS2pHXQ218m1-PjNVK6DLMDPU3btoo32AyvZGWr-G421JbJJGhS-h1galz0CvIJqHw3CPT3f9UoUrPHPFz70Tn1pxTUMQnO5MeDwtdSmVyw0wwtezTYlvscnlbFmlps1KclabYpQSjOgPt-RPozuGs-tPh1xJ4n6dDhBe6uWeMrtxYlh4DZLKXWQ-HWQRicENH_PSl163bJTn8V5xkJXkWUd5kFKfALHaLILASrAzr5TPeapFb9tCjwUWRc7opDiLRUh5Xye3919oIR1vCbrbMgDw24hvgm6CxONlfLKby1BXU7FQsCKl4SWx-VqjEb6W2XGcWFzYRZINM2LA6aECl3dvTlE3ru1dbQw3VReN6gkV4NiQTJt5uf_5enmtj-9XYE2UfpEdgM8XgudXlxuAocFPkX5lwXaAFRlY6_C0IKxnalxxEqfLQHU6D1IRjfkBB8LGqVugcofq5lNarY3AYWF6hmOhDoTUO7TPe4k9qZjpQYZdSfdQjyEfU_nZtx2HS0lKaCWpyGl2IxdiEZ3FKSXy8KFSi2iI2s1YeZq2zFNqeppDrh6h0iGP5UqzDbSG0UKXdmjkXooUGFbo-ZI4h1xuXEkOmsJZUycCawqayhcX5umHUY98k_7Bb1mgr85HRV-MZbEVstgdg8VhrjGoiPtjk0Vl_3aIMBBh3_GbIoFabFVsCpPmY2UBBXUFhewXGeXLqlO0aIne-1NXma3zcmZTjlSSBq3bB6HWXr0FhpvmbalUXMZRuRU3XcTEpy--R1dVbCWy_UZVPilgLLvJOzva0YVSjGqH18dKcYKw5m9Bf6x-pTWJ-ieu51WIAJ-X7naeQ27RbqgLOEDL-L81pExhWlIKk9JuWCciQiG59Ckc0uSNfKVJgMnefkiAt_9nUl6nVY8vvM0NnsA',
--         'https://kinopoisk-ru.clstorage.net/E3Y0U4315/f882714f3/JQb5hfuBU5QkfR6H5sDiqWH4CdAOpqGh5AUpMCYHMTyD-989NpHas7rlbX6SQwnFJE-UnxlPftlyGpwjKZP7jnQa6dHmnQ65USwy44bBk_F82GnlT7uRqLhdQRk_ek1Iup6e9GJJFgNzWY3GateBKfkuhZNkMgvcBl_hRiT_goN1ogDVGx2KcIBFb1uSWN6FNYWtPA5XkZXkeFHEEqdUUHt7UhFGDbLn35HZtsY_MWlhWkWe4EZn2GZP4e74ymbzEc6xncfhMjAkMXP32qyK_fGNJQyCey1lZAh1BGqXmFiT34opeuVu199V6ULme1kRHbINUmw2o2h6p9VmIH5225niYTU7eQZwFHEzX5IQggVVcdGAps-NFVDgBQgml6zsxxt-ASaF7pNPbd1-Jud56dF2uV7ETu-8jnvJTtQ2quuRHkVR8w2i_HjB4--S_D7ZzaWFZHI7OZHYjE0wLrMcwOuHBlVGjXZvt1W9mvYrdWlp6j2udDoTsN5_5d5YOiYbBXZZ4TeJBijQBYcXrnRunSVB1fT6p4llbHTVKE6fAJD3F-ZVionmv8N5NaqqEw2Z-fbBjrTWx3guy1lO8H5ad_HS-QlTiSZonGF7lzIMdvFV8RHIcmeF5fyUfVCyi0xoU3cCVZrl1tfPyUUuMtupPUl6TX4IVvdAZlN5btSW6pf1Ri1N2wEipEDJb2ve0AoFFa05pIY7Uf1UDDWYQvNsrItLEpV-oa7LZ7Fd5kJPFQlxag0CUF7vQBoneUJUDrbDmWalsS8hBoAs9VdXkmBWhWnFaeB6zxF5oIyhLLI7VDBL3541shmqK6sdedLyz_05TVpRatC2G9Bqu01qWG6mT_FarcWv0coolFETJy4I4o3NtXE0jsOFrdiUdczeu4SY74cKzb4dMkcvXWlSBus1OcVqQe4QDkd8rqMVTqh-HhtJYm29v-n6tHSJn5sWzL6h-fntyJ6jSRXUZEkQKntcaMtfRgG-Jebjk8UVzgKrpSkp4vkKZC7v8HqjPaagkoIjdVqx7T_dCkS4qfN3TshiITUl1WR-vzGdoFA9OA4rqDQLw_otZvUu1_-VdZq-tzVpibZp0sSinzzeo31KDB6qd7FOaRVvlTIEuKnvCyKc8llNqX38xm_NqaAcPcAqa8jMZweO2fp9gq9nzbUuBrMxtalmkUpkgs_YWgNNBlgOWl8JGjFl610qiKjZ_18-9Mb1bcll_F4_-TW8lImMPqek6NcTWgWKpabbywGRxrZfeTklXrnaOBYrYCZLafZUbpbDKfbJnZuBxjA0zduf4pDKsdlN_WBaP4GVJPyFNGovEIz_d5q9NmF6T__B7T4yd9WRFTJlAoQquyzuF0Wq9J7SdzUSdem7GQYIeMX7207Q1qmpWYW4Mr9VNTjoDbCaY4Qc90vGTfaFShurRbliSruNue0ywSYYrm_giss1wiwaHpNp1jXNc2HeIIBlY3-2rFapTRVlbNbvqZ2YQA2Aii_ASM-HKr22FVrjb0W1kuZb9d3VFt1-1FYTcAYjgQZQngbHkQ5JRStVfuBAPXOXnvSeremR5ewynwWFvODZsJIvhFzfm8IxBpXCl-eVSc6Op9kVYVKVejxGx3ii0_luQO5m471KRYWrjUJ0zMXjCwpQ9tk1lZXw_s8V5UhQ2RBOZ5Q8awtiyXqFZid_kf1SmseVtRXu9Wq0fkv05gORymx2giMdAm2Nnx22gAwld3eqFIahldGZ2JbzeX3AXBlUin9MyHNPfq1GLUpfy3HtHnK_9SFFUmXujM4LUDJ3mT4kHuJneQKxgTudPoAIfVsPRkjy5XmF6SAyo1EJJFChWEKTkBgDr8aNmrm2A-MNGfYyZ_ndzeLhBnCKh9wWE-WKiHZam_VG9XFv7QIAVN3rrwZAPnmRHb3IkmMB6bQ0xaDmc-yA64P2LfalXhevMVkybivdqUnOhcJM9g80ii9xOjzWXuuR5nHtm1GOkCgJQ-POGBolZZ29UCZDUSk0QJUwPjvI2NPXKrlqZSa_TwXt5sbnSR3t2n1CeI6TDIpTEV6MaoYPkWZ1UbtZCvygsXfztqzG6flVtTDeP9XdJEQVyM43NIx7X9JRHpl2q9eFeXIOo1Uh_VIR2simy2gCA1m-YLqql0nOxbE_WcrksMXvnwokum0lSeXEhjd1ncwApaCeo6S061MKFXKxDt9z9f1qBtvVKSlS5daQjuv8qkudsiwSpv_Jdlld131aoERxF7tCiIZdSTk1ZHoXOWWwsMU0LhcYHGfXmrk-5eKrS2W5Gg4bkclZeuVqVJp7sFJHCQp82gYrBdIViWtRmowMhd83lhiSEWmdqfSio139VECNGDZ7iCiP_1pplmF-A9-RsSpyG72plRIZFli-B3SGT1GGNDYyoz1uMcnD1ZIwhNUzE1r80l3h2aGsgnspMVhc-Tii57AsmxcqoQq9RkO7-bk-xrNV0UnW9UIYsmf44hO9alxKJt9t_i0Vu0EqtHhBi9uWhOaNySF9zHpPsdU8kLFYCg8cSM-fxuEe9cbXn2HNXh5TSclJmtUKkMqDpA4_RT4g7qrfZUL1EWcJXnQ0UedrTiRW9VGFmWxCEz1l5AiNiBIjVOzDa2pRkq2iC6dtaVJ6M4nJCWp1glASo4CiV3n-jMImA4lyPT1ned5IdMX_A6KMkh2R8eGMTleR4agwmVwONxBsV_9qDQ71JqdrRanGametpRVWjUoAdjuMsgcVbrwepuvtem2B8yFGNFiJ1_OSDD5pYYHZJMpHXdXEkJHMLuvMTHOXchUG0bK3T9VBHvJ71Zlpzo1yDP6PhF5Xfe5U7p4P5Y65DcMpgrgU9dubGkDCJbE5lVDyrxkxIAxVWLIboJiDo45h6tGyC8u51d4Oy4E9kVo9-vCWM4C-X0XKpIoGY5HazZFzoSpk2GVfGyaE7vGlXaVwnnNF-XzwESCOb5xo48vmTTLtSrdPBXXqaucJCZmudWbQQhM4GtNY'
--        ),
--        ('Зеленая книга', 'Путешествие итальянца-вышибалы и чернокожего пианиста — комедия-лауреат премии «Оскар» за лучший фильм',
--         2018, '{"США", "Китай"}',
--         '{"Комедия", "Драма"}', 'Inspired by a True Friendship', 'Питер Фаррелли', '---', '---', '---', '---', '---', '---', '---', '---',
--         '{"Вигго Мортенсен", "Махершала Али"}', 'https://avatars.mds.yandex.net/get-kinopoisk-image/1599028/4b27e219-a8a5-4d85-9874-57d6016e0837/300x450',
--         'https://kinopoisk-ru.clstorage.net/E3Y0U4315/f882714f3/JQb5hfuBU5QkfR6H5sDiqWH4CdAOpqGh5AUpIFJuMUGSloNYT4iigraQMDvjNkToCD-VyxgLfoFjQpAqeNa_kzg7ucXygF_wGQhu5t-ku_QU2c0UnnMVQUycuYCmh1BUhv_iLRaJqrvTnVxC6ifhcEGWDVqAjgfY3g9BZlSSplf1Sl3Z190y4Mh9y-faZGYlbfGJsHrf3RU8MFEw5h8klIt7lkES6bq7e0FppppDjcEJSvHiwCpzxHJHtTIsVq7fxcZxZSddQsQIBfsr0ozC_VXNYWxK3xklyBAVWJ5jzNxH9-aFovVSj_OZ3bayy6ERobZp1oAaM4TaDwk61PJSK0nixcVXaaq0BEn7H2bYBgl9kQHsDk8hEfTQiXhmrwhoS1-eyXYBjpPbbZVqgsNdNcV-HXIMdoNUag_lhghqhn-tNjFx0-HG_DxB369K9JpxKfEJdPK3GXEQAEV05heU7Lvv6jV6UYq_S-GZ5nJXNSl1upnW2N5HMJIzbdZMgmIvPb6twSdVPmQkIfOzGugKkSk9DeQiwzX5IFx9DBaHiBCbLwbhsqE-A9tJ3Zoeo5GJnd6Z9piSI8wSrxXOIPLS6wE2bc2nxd7MXN3n_-JManUpLWHUSs-FiaTEtaiyNwSca5OG0YrppiP7bTUq_qcFXdlWEfJUTovY7kPJstSyjsNtwjEZv21GhJQtn98mlFZl6cklrFZn_ZFstCEUhg8U3BfLjsmKldIj1-2tdqrLJdUlojHWcCK34NIPee5UzuaTKT7FWTuZxhAg3VfnulTC6UkFnai2q1G1WASh_LYzJIB_S2Y9GuVOF__hzWYmoxndUSLBGtiaE6SisxFaXGqWi_3Kyc3v4d7EwLHXNxLwmn2xCZkw2i8RpdjwhTS2uxSQ-4veYbYhVk_bEd3echshHSky-SpEessoHovVBnBe2hcB2qkBV-0u-BxFF-_i1EIxfQmh5H6TidlE9MWkplecKB___sk69VZnMxkxQkonGcF5ss3atLbrRPLT0drwkgYvrTq1ed-ZGow0YUcbvhTejVHRcSxK_9m5EEjZNOb3yLjHX3K59lFWJ9-1vWoe4w0ZgeKF4vS261iOv4VKiOqKhzWCZYXrmVaMzEUHe0Z4GvmlTfmAMmeBeaTw3TA61xhAX_9S6RLV9heT4a2aNltZQfFmTfp4pptI2qPtfiTK6p81GjWxd4XeOIBRyxdiyA4teT0hpEbLTV1MQDF4tlsgaM-jxg2qqb4zY-3NVqp7tbkJFpEWwDqPbBp_iXJgfm4HqR41ydcdtjQ4BUOjBuBq7cGB5XjS_40htMQZ1B5rTLQXH_qd5mHiHz9NPRIeZ1EFfTYJ1vh2h0xe08lueA56f3F2tR13AaK8vPUPN5boVrExQQFIhqsJdei81Yw2k0wQM4N-SSoFPm9Xlbne-juVRVn-cQ7QjifU-iu17njqNp-lkuXh36EKvIzlQ3PC0JpdwQGRWH5vCXkYEDX0UqtoDGtPhjW6idbbk-k9xq7DTTnRpkWuEE5_xBID7SZ8TrIfJXaVTceFqmi8_UM31sCGtU2xEcAK59mFRHjJ2JofLERvp5bhsi0mo_v5TaaK7wk1ESadkoTCh1SOl0lOCJK2bzm6xV2ncRpoHCELJ7Z0FhW1zQFkun_dMdhsqZQ6a5Akfy-ubT5p9stf1dVCSk9BHRkSDWZwAmfA8jcNPnAy8mMR0vkxP_kWqFjlE_9CbFIJ0fGpSMLLPSGUhNH0rjsstPsXHi2avYLDq529Ig4rQcEVto3ucAY_7IrbUUo03qYT6XapGUsdGhBULf8jkhyysfEtPbSe40HVfMQJ-FKznDAT61qhFpnmvx8x1ZrypwWF5eL90vBan1wqm1mGqDY-RwHWaUmrjX50rIkfXwr0noFRQSFciq99lbiYRdwmN7BU19cmKf4F2iuvhXWegsOlAXkWQV5gJkv0ZlMBovTCvkeZYkkZaw0KJDxRV3tSzMpdxd3hJCJPSSFsMIlIkpOkrFfjXrXGBaZLyzXJRmbDJQXFNknaDK7zwHYrtX44XnZP-Zo1nZ8dDqTEoVuHBmRCpS2pHXQ218m1-PjNVK6DLMDPU3btoo32AyvZGWr-G421JbJJGhS-h1galz0CvIJqHw3CPT3f9UoUrPHPFz70Tn1pxTUMQnO5MeDwtdSmVyw0wwtezTYlvscnlbFmlps1KclabYpQSjOgPt-RPozuGs-tPh1xJ4n6dDhBe6uWeMrtxYlh4DZLKXWQ-HWQRicENH_PSl163bJTn8V5xkJXkWUd5kFKfALHaLILASrAzr5TPeapFb9tCjwUWRc7opDiLRUh5Xye3919oIR1vCbrbMgDw24hvgm6CxONlfLKby1BXU7FQsCKl4SWx-VqjEb6W2XGcWFzYRZINM2LA6aECl3dvTlE3ru1dbQw3VReN6gkV4NiQTJt5uf_5enmtj-9XYE2UfpEdgM8XgudXlxuAocFPkX5lwXaAFRlY6_C0IKxnalxxEqfLQHU6D1IRjfkBB8LGqVugcofq5lNarY3AYWF6hmOhDoTUO7TPe4k9qZjpQYZdSfdQjyEfU_nZtx2HS0lKaCWpyGl2IxdiEZ3FKSXy8KFSi2iI2s1YeZq2zFNqeppDrh6h0iGP5UqzDbSG0UKXdmjkXooUGFbo-ZI4h1xuXEkOmsJZUycCawqayhcX5umHUY98k_7Bb1mgr85HRV-MZbEVstgdg8VhrjGoiPtjk0Vl_3aIMBBh3_GbIoFabFVsCpPmY2UBBXUFhewXGeXLqlO0aIne-1NXma3zcmZTjlSSBq3bB6HWXr0FhpvmbalUXMZRuRU3XcTEpy--R1dVbCWy_UZVPilgLLvJOzva0YVSjGqH18dKcYKw5m9Bf6x-pTWJ-ieu51WIAJ-X7naeQ27RbqgLOEDL-L81pExhWlIKk9JuWCciQiG59Ckc0uSNfKVJgMnefkiAt_9nUl6nVY8vvM0NnsA',
--         'https://kinopoisk-ru.clstorage.net/E3Y0U4315/f882714f3/JQb5hfuBU5QkfR6H5sDiqWH4CdAOpqGh5AUpMCYHMTyD-989NpHas7rlbX6SQwnFJE-ckwl_brUHW8Q7NMfvnnwbrcC3xFPxXTFG5sLAy9lwzGCkAu-E7flpWF1CXlVYJosy9aoxYntLGdXmBstB0RxK9epkIm_YHlPwWiCGRjYhuilBI1k2kPR9z7u-hMYloY2R5DLjsfEoRIVEWp-wlMOvdlFOnaovO7U9VkZTNRERzoGGYEJ_2LaPxb5Q4iqHaWbV-WP9QoxYNTvvxkDOrZEBvVjCY8HV6Dy1iFJ3FEz7k56Nfp1uH8-VeT4-L91Z3ULxQtBel-w-V3GueGoGV8GaTc0jzQLM8H2H5z7kMlkdJQn4slcppeRwtbzmI9C408_-DToNVivzVeUexuMZ7dHqiQ4EqkvwFqM5ckhi-nOlUjlpr6GyHEB9a1vifOYN-fH9TDbfRe3ceJEMyg9MwIev9pXG9W5LF4UpEkZbhWkhWv3yCPpP3IYvNf649pJvFZa9zXsJdni4QeMLppQCXWl5YfzCa711xBi9EJoT3CCHY_IFFoFCwyfZEWq2y5mVAZoRJsAK-2AWh3GC1AI2z_3yve07RRKEON2bE8rksplV8aHwQvtd3bzkqVxit7zEh3OeNX6N8rOjQdnOEnsVGfEmkRb4QmNANqOZMjQGohu5ejXp95m6kMQxR28-pO6xOQX9JFpTxZV0FNF8pm-A1EeX2k1iJYqrazFNciZDBVmNfpkO-D4XQBojAW5gaoKTRY4VzdP1hqj4ffczvtiG4X35CWTep0UBwOQZRDqvFFjnW2JJgukmj1-BzZoWfzUF5f5x-mhOi3QyL2F-7AK-mzEO5QF7TSLsiMGfh7Z89vmpDQXwCt9d1SCImZSSC0zMH1dm0e5tZp_fdelSFvcFFWE-yabEipMsFt9xxri6hltJHt0x5636YDT5W9uaSLplVR1lPLLTren8fFlMYi-UgNNXXgVK0f7jQ3GpwgYbja2FSukOSF6TBP7XnVqAhr6HGZ7pwRdh2gzYoV8HGoRmXfn9eUQ6p5md1FgJuD7vCDz_j47Nfr2ugxfNtVJGu9k9XeplfoT6k0QSexFy1EKqX-HOoflXYdoQpM0Ll2L86vVhRam4DqfVnSx8SdjGg8xICxMGYQYl9kOjdbFWmpsJxcVKRS5gfjN0Xi8Bgvz6_geRSmnh23GqAPDRY6PO3IrtYd35jJK7XSlgaIW04jPYnNdj3kVyiTpnS8VdHhYXMe1VFtHK2AJ7UK4jYU5g2hL_aTq1DWPtviQwDQevimgOdf3Z-fQyIzUl2DwNAIYbvFxv3xqZ5r36G7NBdbK-J10xjartWpTKJ3zyg5EK1Mb2Qx0aLc1bobYEdKFHs5IYGg0lsXkgkj8hrVzMQZQWE4AAnx_-qbLpfk_vObnqlt9dlak2aY5YrvsMmlsVxjCaMgM50lUVc1kWnNBZOzOS_Fbt8VUp3Dqfia1s3A3QQitM7G9fbrlKLX5DH5VZkvLneYnx-pHyyCITuF4nkd5kYup_sYphtbOZTow4cWP7lljSbXGxWXAiuyl5XMQNlFY7UATj7-4hPqWuv0P9pb46Uz3B9RKBJsCG48A2N-G-QE6uc3EKuYknFbYcpOXHk-KE1h1tfQlgQk-ZefwYRYQ2j8CkG5P-hY49qgvf6cXymieBoeWauapMwjOokht5WoDu5lt5Pil909VWiNhFg-OaJJIRRRU1DNrHlbm43F1cwpeEuH-vVqn2iUobkwG9kg53PTFhognq6BZHoGZTETrEiuaHdZqp9dPRDqSgqd-X3sjGYb2xZSSuI5kBtBSxgBLnZABfc8JVqqE273tBZZ7y_421iV5NZmQyI9zS_3mCOAaiw4XO2clTja4UAOnXW0IgXjVVEaV0TrP9ZUywUfyKD0gw_x_evb7tCq-_HSm6hnuh0U1iMe6Mrh9IYkvZhkhiAkcZOmVFw_F6vEwhj38e1N41zaWFJI4ziTXcaBnY0jcc7GuDHsUWDT4ba7XlLjLftSnNVklytK5jKAb7ZV6sYoJDpRptwa95wohcWTuj0kgWPa1d-aB6I421JJgVJIaflBSD9-KVApW-j_99oTIOzz1FVeZhKtAmM2DmF7VyNLoq80WebQG3abYQMOWz31aUCm1ZBfEAOsvJBUzIgbS-D5jMx5vK7XYxzgvnddmyBhs9sVm-SQpEjnuk6lsdflw6km-pdkmR850C6BStH-Nm-Hq9-fnRTMK3eWXYeDUIFoMcXGvXngECCV5Pl30Z9uZrFbHlel2aCHZ3MFIL1d6I9jYjfcplUd_V9iCYeY_3KtjeIWkhZShaU4kt9GBZmCJrNJy7fxqdqp2qR6cBGdqGp31NmXZ55syif2jeQznqAM6KBz1i4VljXabMvLVrt2ZQmikxAb1cll-VWdT0xaAmf9zsc-PGper5wk-ztbEy_nu5oc02dYZAxiOEMitF_nyeGhvhGnXh56EydHR5E4O2eGL1UfmJxHI7WRG0XC0MQitUADP3jiV-3Vo7021RLuZ79YGFvg1iHCoPfGZX4XJ8lqbD5cY9lSftIhjEobMzzuDGEfHB1UjC48EtZEQBROYnoKyDe9ZBouVWn98JMe7mOwUhDX7VQjiGZ0Cm-83-oHqWC8nGTRUbrbYArE0b9yYgsmkRzZHkRq_5ObBYFQBmszSs3-eOxQ4pfl9LGWXKiic52cUusdo0ljcsNssRfkgenlt1UhWNZ4H6KFx9m1tS0MJRuUmBKHLDWTEgeMncRpdctMfvqlEeDe63k4F5srZbodn9IjluPHpnRLYj4UasFmqP-WIdSevNhiQ09denHgB6Hc1xaWyWJ8X1tOQ5sJJnaEizA6pRoomCI1N9yeYSozVpdd5R0jiab3yS04XewGI--2XSleE3ARagtMkTi8oUHi3tHbUwXns5sczYTYxiBwAgn9uWqR4NPoNnGeVuJqvBIen-hfKAPuNg6rdU'
--        ),
--        ('Властелин колец: Братство кольца', 'Фродо Бэггинс отправляется спасать Средиземье. Первая часть культовой фэнтези-трилогии Питера Джексона',
--         2001, '{"Новая Зеландия", "США"}',
--         '{"Фэнтези", "Приключения"}', 'Power can be held in the smallest of things...', 'Питер Джексон', '---', '---', '---', '---', '---', '---', '---', '---',
--         '{"Элайджа Вуд", "Иэн Маккеллен"}', 'https://avatars.mds.yandex.net/get-kinopoisk-image/1629390/1d36b3f8-3379-4801-9606-c330eed60a01/300x450',
--         'https://kinopoisk-ru.clstorage.net/E3Y0U4315/f882714f3/JQb5hfuBU5QkfR6H5sDiqWH4CdAOpqGh5AUpIFJuMUGejq9Yb4iigraQMDvjNkDQACLMqxVONql6E8gnDZqjqzgG-J32qRKhdSxu5t-ku_QU2c0UnnMVQUycuYCmh1BUhv_iLRaJqrvTnVxC6ifhcEGWDVqAjgfY3g9BZlSSplf1Sl3Z190y4Mh9y-faZGYlbfGJsHrf3RU8MFEw5h8klIt7lkES6bq7e0FppppDjcEJSvHiwCpzxHJHtTIsVq7fxcZxZSddQsQIBfsr0ozC_VXNYWxK3xklyBAVWJ5jzNxH9-aFovVSj_OZ3bayy6ERobZp1oAaM4TaDwk61PJSK0nixcVXaaq0BEn7H2bYBgl9kQHsDk8hEfTQiXhmrwhoS1-eyXYBjpPbbZVqgsNdNcV-HXIMdoNUag_lhghqhn-tNjFx0-HG_DxB369K9JpxKfEJdPK3GXEQAEV05heU7Lvv6jV6UYq_S-GZ5nJXNSl1upnW2N5HMJIzbdZMgmIvPb6twSdVPmQkIfOzGugKkSk9DeQiwzX5IFx9DBaHiBCbLwbhsqE-A9tJ3Zoeo5GJnd6Z9piSI8wSrxXOIPLS6wE2bc2nxd7MXN3n_-JManUpLWHUSs-FiaTEtaiyNwSca5OG0YrppiP7bTUq_qcFXdlWEfJUTovY7kPJstSyjsNtwjEZv21GhJQtn98mlFZl6cklrFZn_ZFstCEUhg8U3BfLjsmKldIj1-2tdqrLJdUlojHWcCK34NIPee5UzuaTKT7FWTuZxhAg3VfnulTC6UkFnai2q1G1WASh_LYzJIB_S2Y9GuVOF__hzWYmoxndUSLBGtiaE6SisxFaXGqWi_3Kyc3v4d7EwLHXNxLwmn2xCZkw2i8RpdjwhTS2uxSQ-4veYbYhVk_bEd3echshHSky-SpEessoHovVBnBe2hcB2qkBV-0u-BxFF-_i1EIxfQmh5H6TidlE9MWkplecKB___sk69VZnMxkxQkonGcF5ss3atLbrRPLT0drwkgYvrTq1ed-ZGow0YUcbvhTejVHRcSxK_9m5EEjZNOb3yLjHX3K59lFWJ9-1vWoe4w0ZgeKF4vS261iOv4VKiOqKhzWCZYXrmVaMzEUHe0Z4GvmlTfmAMmeBeaTw3TA61xhAX_9S6RLV9heT4a2aNltZQfFmTfp4pptI2qPtfiTK6p81GjWxd4XeOIBRyxdiyA4teT0hpEbLTV1MQDF4tlsgaM-jxg2qqb4zY-3NVqp7tbkJFpEWwDqPbBp_iXJgfm4HqR41ydcdtjQ4BUOjBuBq7cGB5XjS_40htMQZ1B5rTLQXH_qd5mHiHz9NPRIeZ1EFfTYJ1vh2h0xe08lueA56f3F2tR13AaK8vPUPN5boVrExQQFIhqsJdei81Yw2k0wQM4N-SSoFPm9Xlbne-juVRVn-cQ7QjifU-iu17njqNp-lkuXh36EKvIzlQ3PC0JpdwQGRWH5vCXkYEDX0UqtoDGtPhjW6idbbk-k9xq7DTTnRpkWuEE5_xBID7SZ8TrIfJXaVTceFqmi8_UM31sCGtU2xEcAK59mFRHjJ2JofLERvp5bhsi0mo_v5TaaK7wk1ESadkoTCh1SOl0lOCJK2bzm6xV2ncRpoHCELJ7Z0FhW1zQFkun_dMdhsqZQ6a5Akfy-ubT5p9stf1dVCSk9BHRkSDWZwAmfA8jcNPnAy8mMR0vkxP_kWqFjlE_9CbFIJ0fGpSMLLPSGUhNH0rjsstPsXHi2avYLDq529Ig4rQcEVto3ucAY_7IrbUUo03qYT6XapGUsdGhBULf8jkhyysfEtPbSe40HVfMQJ-FKznDAT61qhFpnmvx8x1ZrypwWF5eL90vBan1wqm1mGqDY-RwHWaUmrjX50rIkfXwr0noFRQSFciq99lbiYRdwmN7BU19cmKf4F2iuvhXWegsOlAXkWQV5gJkv0ZlMBovTCvkeZYkkZaw0KJDxRV3tSzMpdxd3hJCJPSSFsMIlIkpOkrFfjXrXGBaZLyzXJRmbDJQXFNknaDK7zwHYrtX44XnZP-Zo1nZ8dDqTEoVuHBmRCpS2pHXQ218m1-PjNVK6DLMDPU3btoo32AyvZGWr-G421JbJJGhS-h1galz0CvIJqHw3CPT3f9UoUrPHPFz70Tn1pxTUMQnO5MeDwtdSmVyw0wwtezTYlvscnlbFmlps1KclabYpQSjOgPt-RPozuGs-tPh1xJ4n6dDhBe6uWeMrtxYlh4DZLKXWQ-HWQRicENH_PSl163bJTn8V5xkJXkWUd5kFKfALHaLILASrAzr5TPeapFb9tCjwUWRc7opDiLRUh5Xye3919oIR1vCbrbMgDw24hvgm6CxONlfLKby1BXU7FQsCKl4SWx-VqjEb6W2XGcWFzYRZINM2LA6aECl3dvTlE3ru1dbQw3VReN6gkV4NiQTJt5uf_5enmtj-9XYE2UfpEdgM8XgudXlxuAocFPkX5lwXaAFRlY6_C0IKxnalxxEqfLQHU6D1IRjfkBB8LGqVugcofq5lNarY3AYWF6hmOhDoTUO7TPe4k9qZjpQYZdSfdQjyEfU_nZtx2HS0lKaCWpyGl2IxdiEZ3FKSXy8KFSi2iI2s1YeZq2zFNqeppDrh6h0iGP5UqzDbSG0UKXdmjkXooUGFbo-ZI4h1xuXEkOmsJZUycCawqayhcX5umHUY98k_7Bb1mgr85HRV-MZbEVstgdg8VhrjGoiPtjk0Vl_3aIMBBh3_GbIoFabFVsCpPmY2UBBXUFhewXGeXLqlO0aIne-1NXma3zcmZTjlSSBq3bB6HWXr0FhpvmbalUXMZRuRU3XcTEpy--R1dVbCWy_UZVPilgLLvJOzva0YVSjGqH18dKcYKw5m9Bf6x-pTWJ-ieu51WIAJ-X7naeQ27RbqgLOEDL-L81pExhWlIKk9JuWCciQiG59Ckc0uSNfKVJgMnefkiAt_9nUl6nVY8vvM0NnsA',
--         'https://kinopoisk-ru.clstorage.net/E3Y0U4315/f882714f3/JQb5hfuBU5QkfR6H5sDiqWH4CdAOpqGh5AUpMCYHMTyD-989NpHas7rlbX6SQwnFJE-ckwl_brUHVpgmfYvu3kAHncX2mRaoBQ1a4tOhm-l8wGysFubE8KVpWF1CXlVYJosy9aoxYntLGdXmBstB0RxK9epkIm_YHlPwWiCGRjYhuilBI1k2kPR9z7u-hMYloY2R5DLjsfEoRIVEWp-wlMOvdlFOnaovO7U9VkZTNRERzoGGYEJ_2LaPxb5Q4iqHaWbV-WP9QoxYNTvvxkDOrZEBvVjCY8HV6Dy1iFJ3FEz7k56Nfp1uH8-VeT4-L91Z3ULxQtBel-w-V3GueGoGV8GaTc0jzQLM8H2H5z7kMlkdJQn4slcppeRwtbzmI9C408_-DToNVivzVeUexuMZ7dHqiQ4EqkvwFqM5ckhi-nOlUjlpr6GyHEB9a1vifOYN-fH9TDbfRe3ceJEMyg9MwIev9pXG9W5LF4UpEkZbhWkhWv3yCPpP3IYvNf649pJvFZa9zXsJdni4QeMLppQCXWl5YfzCa711xBi9EJoT3CCHY_IFFoFCwyfZEWq2y5mVAZoRJsAK-2AWh3GC1AI2z_3yve07RRKEON2bE8rksplV8aHwQvtd3bzkqVxit7zEh3OeNX6N8rOjQdnOEnsVGfEmkRb4QmNANqOZMjQGohu5ejXp95m6kMQxR28-pO6xOQX9JFpTxZV0FNF8pm-A1EeX2k1iJYqrazFNciZDBVmNfpkO-D4XQBojAW5gaoKTRY4VzdP1hqj4ffczvtiG4X35CWTep0UBwOQZRDqvFFjnW2JJgukmj1-BzZoWfzUF5f5x-mhOi3QyL2F-7AK-mzEO5QF7TSLsiMGfh7Z89vmpDQXwCt9d1SCImZSSC0zMH1dm0e5tZp_fdelSFvcFFWE-yabEipMsFt9xxri6hltJHt0x5636YDT5W9uaSLplVR1lPLLTren8fFlMYi-UgNNXXgVK0f7jQ3GpwgYbja2FSukOSF6TBP7XnVqAhr6HGZ7pwRdh2gzYoV8HGoRmXfn9eUQ6p5md1FgJuD7vCDz_j47Nfr2ugxfNtVJGu9k9XeplfoT6k0QSexFy1EKqX-HOoflXYdoQpM0Ll2L86vVhRam4DqfVnSx8SdjGg8xICxMGYQYl9kOjdbFWmpsJxcVKRS5gfjN0Xi8Bgvz6_geRSmnh23GqAPDRY6PO3IrtYd35jJK7XSlgaIW04jPYnNdj3kVyiTpnS8VdHhYXMe1VFtHK2AJ7UK4jYU5g2hL_aTq1DWPtviQwDQevimgOdf3Z-fQyIzUl2DwNAIYbvFxv3xqZ5r36G7NBdbK-J10xjartWpTKJ3zyg5EK1Mb2Qx0aLc1bobYEdKFHs5IYGg0lsXkgkj8hrVzMQZQWE4AAnx_-qbLpfk_vObnqlt9dlak2aY5YrvsMmlsVxjCaMgM50lUVc1kWnNBZOzOS_Fbt8VUp3Dqfia1s3A3QQitM7G9fbrlKLX5DH5VZkvLneYnx-pHyyCITuF4nkd5kYup_sYphtbOZTow4cWP7lljSbXGxWXAiuyl5XMQNlFY7UATj7-4hPqWuv0P9pb46Uz3B9RKBJsCG48A2N-G-QE6uc3EKuYknFbYcpOXHk-KE1h1tfQlgQk-ZefwYRYQ2j8CkG5P-hY49qgvf6cXymieBoeWauapMwjOokht5WoDu5lt5Pil909VWiNhFg-OaJJIRRRU1DNrHlbm43F1cwpeEuH-vVqn2iUobkwG9kg53PTFhognq6BZHoGZTETrEiuaHdZqp9dPRDqSgqd-X3sjGYb2xZSSuI5kBtBSxgBLnZABfc8JVqqE273tBZZ7y_421iV5NZmQyI9zS_3mCOAaiw4XO2clTja4UAOnXW0IgXjVVEaV0TrP9ZUywUfyKD0gw_x_evb7tCq-_HSm6hnuh0U1iMe6Mrh9IYkvZhkhiAkcZOmVFw_F6vEwhj38e1N41zaWFJI4ziTXcaBnY0jcc7GuDHsUWDT4ba7XlLjLftSnNVklytK5jKAb7ZV6sYoJDpRptwa95wohcWTuj0kgWPa1d-aB6I421JJgVJIaflBSD9-KVApW-j_99oTIOzz1FVeZhKtAmM2DmF7VyNLoq80WebQG3abYQMOWz31aUCm1ZBfEAOsvJBUzIgbS-D5jMx5vK7XYxzgvnddmyBhs9sVm-SQpEjnuk6lsdflw6km-pdkmR850C6BStH-Nm-Hq9-fnRTMK3eWXYeDUIFoMcXGvXngECCV5Pl30Z9uZrFbHlel2aCHZ3MFIL1d6I9jYjfcplUd_V9iCYeY_3KtjeIWkhZShaU4kt9GBZmCJrNJy7fxqdqp2qR6cBGdqGp31NmXZ55syif2jeQznqAM6KBz1i4VljXabMvLVrt2ZQmikxAb1cll-VWdT0xaAmf9zsc-PGper5wk-ztbEy_nu5oc02dYZAxiOEMitF_nyeGhvhGnXh56EydHR5E4O2eGL1UfmJxHI7WRG0XC0MQitUADP3jiV-3Vo7021RLuZ79YGFvg1iHCoPfGZX4XJ8lqbD5cY9lSftIhjEobMzzuDGEfHB1UjC48EtZEQBROYnoKyDe9ZBouVWn98JMe7mOwUhDX7VQjiGZ0Cm-83-oHqWC8nGTRUbrbYArE0b9yYgsmkRzZHkRq_5ObBYFQBmszSs3-eOxQ4pfl9LGWXKiic52cUusdo0ljcsNssRfkgenlt1UhWNZ4H6KFx9m1tS0MJRuUmBKHLDWTEgeMncRpdctMfvqlEeDe63k4F5srZbodn9IjluPHpnRLYj4UasFmqP-WIdSevNhiQ09denHgB6Hc1xaWyWJ8X1tOQ5sJJnaEizA6pRoomCI1N9yeYSozVpdd5R0jiab3yS04XewGI--2XSleE3ARagtMkTi8oUHi3tHbUwXns5sczYTYxiBwAgn9uWqR4NPoNnGeVuJqvBIen-hfKAPuNg6rdU'
--        ),
--        ('Побег из Шоушенка', 'Выдающаяся драма о силе таланта, важности дружбы, стремлении к свободе и Рите Хэйворт',
--         1994, '{"США"}',
--         '{"Драма"}', 'Страх - это кандалы. Надежда - это свобода', 'Фрэнк Дарабонт', '---', '---', '---', '---', '---', '---', '---', '---',
--         '{"Тим Роббинс", "Морган Фриман"}', 'https://avatars.mds.yandex.net/get-kinopoisk-image/1599028/0b76b2a2-d1c7-4f04-a284-80ff7bb709a4/300x450',
--         'https://kinopoisk-ru.clstorage.net/E3Y0U4315/f882714f3/JQb5hfuBU5QkfR6H5sDiqWH4CdAOpqGh5AUpIFJuMU2CmodEf_jXz_KQMDvjMljVSDOV3xgHeqljSow_ONffizAXpLCaqQvkFQlKgseJi-kUxGSw5gsZOXioodSmo6iwh5uTMQKR0ru37VU2j0tV2b2P2SqIivN0EiM5cvTWhos9winBX8U2uDwhi68anIIZyQUpgKKv_ZWw_NF4ThPoKPNbnrV2_dbbp-396rqvJb3RPpH2dDKz0GY_lToAgv5PNUoZTXN5xjhMBUvXKlCK8W3dEbxKc82VdMwlWAp7kFQbE1I5Bjlmx0_ZdTIOvw01_e45CuwG8-Amfz1yvIoG68m-lWnH2bYMpHVHmypkPqWpKTngKvOJBUz4GZiWW2iY36dekX51sjOTxV3GRmM9PQHKXcKYon-Mlq-NclA22nMd6nG9M20yhMg9f5MO1BKJNVFtgCJrdf10mP1IWlfoIEMjriEKib5jl-nNSkrvzalp1u0GHAarJFLLdU7YZp6b-brhNa_dxjAwpWfzIshClaWxbUwm-6WJWBDNFGIvGLBf347h5l12kyNVXeIOk6FdzXYFYhwm62g2N_XSoH7y60l-3b1v0Uag0A0fDzaEujHFVW1cSsvNhehgSYyqi7wA01N-XWZtTtu7dX3G5iNBWVmiQeqUIie0niMJPnwCBqsVVrFJMwVeCEhF1_9OpH7p-UWtuA6z0S2QeIH8PjeIOMMTAgVudU6nz3VRRn5_FTV5Kr0etAYD2KIbNXLMXobXfQb1tcdF2vzI0WMPhpziKW3JDXS2tzHhPFy1TL7fuATzT2qFhoHe11NBeUoeb5ldRSLJnkTKq2AGX0XOpOqOcw0eIUHL0Q6E0AWDYwZMSo01XfV4si9dZXxMNbiaF7iMw1_uRT7dchNLGV26DtfN5X3isY58-jeA3tP59mC2okdBgt1Rqx22iCA5X5fGlLqp7RE5eIr7-dnkMKm82oeoYEvnCjEedf7HSzG1suJL9dlFPuEOSArHTP6_Fa5kaiKLnbpxsbdlPvwUTXezlmDmaXGtFaBaM821tFD9AMYX6MAfd9KRkgUyY0txWR5uY6EdUeYZXgAyh0z-o2nCMPpa8xES6QlnmQr8WE2Pl9YAHgW12eE80p-1LeyQSbjCEzTgz49KMbJV1ufrQRVKfpOJpQW-adrIKgtcjrM93ljO9tNxCumRN62W4ND5w4MabDq1oQ09TAq7wYEgtKEILlu4bPen2m0msW6bo2XlRh5fFYXpRpGqFMazwJqX_QI8wrJn9ZJ1lTfVNni49XvXkthencXNhfDOZ1W14MhZjAb3EFybewLRGiEiU_9JuebuG6GZDfrliowGi4ySt7mufN6qF-Hqrf23AZZkrH3_J95MzpX5kXUwKlcB4WScBfTKrzikm98mTZ717jcjOdE-atdFxcm6wUL03qN0Mi8dVgBeqvOtCnkZ5_0-xAR9zzeSCJqtNX2FcLpH-SVkkPVYKtdcnL_DfoFmiX67y40VQu7PET0RxkkawH5jtGo_9X5Ylq5XKYr5_ZdRJuCkqf8vkkyOvSmVCcA6342ttGypMNb7lCj7i3ppdl12Hzv1fVKerzURVcqJmhhC9ziSr2nq_P7aiy365THHQUYUFKlf89pc7gm5NfG8Kns9NbDYNSS2tzRcR-tq4U7R-lvrndl-Bkv1sR3iga6ItgP4cjsVSriOoitp9s1Z-y3enBhpGzfChBoR_SmVgIJXRYFQyHnMztegDPt77tn-kV6Pn5UtNm4rsdUdPo0KCD4D_CoXbabk-ubHPYY1_asFqngU0Rf_LljKYR2RtVwWqxmpLDyRjBbbXIRL_wYluh3Sq_vpmZoGk01ZWXp9XngCg6CKp83m7DZ6L6XS3V1rVUrocLXvW84kUokxoRUwCkMN5RB8VdBa_ygAZ5vCGcaVOjfHfSkuppc9Pfn-4arEjhPcXg-BLrQSJtsl0kXpSwWKaATlf4OGAAqxZX2BrMo7pQUkyIF4lmucpHNjQi2-CQI3ux1NnhpP2T15-l2KzAp_VOY7kVYAzupH7dolETeBfngAZYdzivxeGe2Fadg2a7GdpFwVsNJ3oLT7D9qdllFmv-tVrXLKY0Hl0Uq9DszKZ0SSo_3qiLJum_GK0Uk_IT6QRNXvIx5sZonhXS20HhPFOdTYDbiq96hg-_vWxb5x8hejkaE-Ym8pZWnWUeboWiOwJlvZoiSOXveBWnG1H23G7PS1e5Oq0M4FZc2B-Er_sQFEnH2warNIENP7agGq4b7vrwUZbqrP_anNmoVaxJoP-NKTVXa0mhLXJcbhbasJXggE_VeLxkD67U0NUVDOYxmVsJRNzGqfKNy7BxYNjp16O6ddlSZG-3WRcb7F8kCSs3CCf3G6UNpeX2HOuU1zfZIEGIl3H1p4_vmlfZnMEltZ8dicWXjCd1AAf-tCTYL99l_7sXlOOu8JweGiGYrUKjeMFse5dijujneZEtm1R-V2YNTBF7ey1JqtLZHZ2FrbzdVA6DmgImtIADPLCsX6Gaqz10ktMp5jCcldeh1WnF73wAarCa6IXvbvPfZ5jRtpxrhM_cevnpw-odk9aVQCvxHtTEw1xEKrSEDDa4IFIjmOH7917Z6y79UlbbIxVuzey4CSs2FCIJoeL0mOmYFfxUL0dOkTs4rYvjVNPTXIWju9IWSModQWjyRc_5NKVUahgg_vGX2ubm89QWXijcK0Rres3puRcqA2at85tjEFTwl2mNThg5NWBJ4RJSUtwH6vrQX0ZHlMCvcYIGeTclnOFYrjv3H9Rp5X2UmRNgHyvII74KKX-frsyiYPgfpFPadNknxIJRcPpmhK4RHZWSx-rxGBmPC5sLqjvNjzI_qlpqmOA7dJ2bb6z7U9xUKdQjQq5ywyE3nGKObyG-XKZVF7EVogtGFvM9JUuoF5sXX0QletBSRQjdSWK4jQB2tmhXKJNqc7VaHSKiu9IaFi0cYYhk9E5s_Q',
--         'https://kinopoisk-ru.clstorage.net/E3Y0U4315/f882714f3/JQb5hfuBU5QkfR6H5sDiqWH4CdAOpqGh5AUpMCYHMTyD-989NpHas7rlbX6SQwnFJE-Ugz1fftl3WowjLbfiwzVO8cCjxEP5TTFe2sOYy-wQ2FSpQ5LA6M0RXFxjZlz1lzsyjao9Fjs_de3eGi9B0Hle_fZgXhPAdjLlKjwuR8vFimkJb-EiSADxV4PC2EZtYTk9wA7b0WXgzEVEooOQELN_lmkG9UJLk51dnoJbgd39KpHyAE4TaKoHAVpYQvaDGXbRScuVPuRIBQP7BtDOXe0VgTCOq_WlmPyJTEonSCiPl0pZBjFyv7PZNeb-s8kRcVpVQhymJ-BysxFy0G4mK-Xu5Qn71X5MALkLA6IsOtHJoSFAukOFqdT8vfge47wA0_fKHZYJRoNzRRUeMnd9HdkiGZboejvIhvvNQtiSAk8tmkGFl2Wu_ABVt986-G41HVWVxDIvzZHc2A3UMn_EVLP_UuFuMSZno4kZnorr-e1pVuWauH4XWAr3QbJM-h7_6R7lUT-hygQ83eeb0hw-pZXJJTCG11WJvPQRhC7vJFR_-8IxGh2uV_-xYW4a9wXNqboxUkjKq8iisz3euF6-F40exRFzxTaEoKX_96Ks-pkdCSmwFjf98UDgXXyKj8BUb5fyWRat3tNnecXKqnuJPRU6AWoAUovohluNPrzKalMFlsHdr20ieEx5gwPi8NL16VX9qL6vtTmwmH24UrPQlIvTikW-1cYbF-15_pJryUFNMhlqfCaLxAbD0WrQ6uKv8bbl-cNRGkQAyd-DnpiCsRWhvSxKLyGNQFBFJJInXDRHa46lcnniL6dtkc6uW5Upzdrt-gy6v-wKo8HmuNbq23FGKVF79V40vKFrizrommXhrSn4Mjf1bSzQlYw2f8jMS28WyfY58q9TSVnOJmuFrQ1isVbIoufI-rN5sgDuKqNhfhnNmy3SiIRlN6cOpAaZ8c3lQD7HybHYEE18EqeEAEtXwm1KoY4zVwnJ3srjPUl5QhnaHKLPIPJf5Yo81vbz4UrpPVcNvmTcYesnwng-NRHRnchK872Z_EC5INI7OCyThwpZJvHuZ-sVWZ5qt62R2c5pFriij8xe083e-MIuC7EC0X1XDaIYsDV7X7r0lq2pAWH8Sr-9YdgA2di-_0zYDw-mIb6pLtNTEV1CSmdVCXnuOfI8Ar-ACsM99kCWdns1ysnxR32yTKxdT_OalI6tMVFVYFY3CS3MzLX8DuuYBH_XglUSZQo74_0VzsZffZklet1KQEqbcAaj8Wpgeo6DRRYlSdtploxwOUO3LhAWMTVRLcDOXwWVmEQBmCaPWLzDE17BJqV2w2fVuWb2M6FBmUZNBogWtyymU7XefJ4y92WO5XGXYbbI3Hlfr14Ebuld0flg0kuNEWgIlQguswRMA_dulXIhIp8fGeFODjMFZQXCmcrsysdEftd5OiBactOt9j1Zb8EubCQF36-6SI49uYEFyHLjjSF4RNFcFn_ovENnfm22IS5vs_mZKjYXGT3JOuVaYCJzgAJTYW7Ygg5b9cKdma-ZPoQMXRerHswOvV3xqdBWQ1kRYESVSAZjADDz5-YZPvHSM9sFteKCU1E5ISoxUsTSC-gSIwFK9MYCm3UaoQ0jYa4YmPl_38LIfqGRobmwovNZsbwMhSiy86DIj_dCqab1Zq_PZflC9u8xKakSvd6AAmNMPrvlilSOKpNBilX544E6ZDi9D6dijHKJ-Z3VKCr_mfV4FF3cqre8rLNfbtESFXbjJx2Z1qZToa2Rov16VHZruHbThc4wjvaf5Qrd-efZFhzU4XvjjtgCcV3N_VzO8yH5sPiBDNpXBIxvy5KNOmmCC2fFlSou4yVFbeZx9nASFwzauz0yvMqyb7F64Xm7eaa8lOm3f2ZAVpn9Da28XpdFARQY_ZQyezQsA9d6mXZVws87ibFeqs9BgVGa-R7sLoO8bhs5QthqNvNFxm3px60O8FyxkyOSwFYBSS39fN7jFZHMUNnMCi_ouJ8XAjGWYXYbk0Ul6g7buQFl4mUm7FLj2N6n4abY6jJPZc7phU8VOuAkBU_vDgheYbFReYjO55VpPFwlmKKnEFDr61IlDuHij1sBOdYeU9WZ1co9QmQCqzgyd80-AEKCr-HOKZ1fYaKMmI0za9IUDpXpWdnIJqMlAWzItaAyq8gUh8MqUaqRZpdTebneylMhlY3iHdbMSm80ft_BVoD6HkMJ6rnZq9VaqNAhD1u-ZN41FXmVMFoTRZXcfAkIvi9YuMuXxiWSASLnW7n9Prp7ISlJ9o2aNEb7jC4XYYJMXlKXtcZ59eMhkiQEsRsXnsBCpc3N8ai-4w25xBCZPFYHmGhjE1qNBvUq1ye50V52E91VRdLxXuBOowBm-1UKdOJ21x1CcUlrcX4AyFVbWxaESv3tFYVksv95mVCMoThC7-ig_89izWKdIsOTETkmqtcxAQXekdKEEk_sDodBdiRyagtl1snNl-XGyAQtb4s-fJadFSEdgNYzMfn4ZA1cFmcE4OuH4llGBVajS_ElPqqbEUmNpnWOaD63uHIjzXYszrIPuZ69Ddv1qnjcjd_zpthyPS19kTAOqw0p4EhF-BqTqFBn34aFfgnyry-R5T7qa7HBTX5VqsRWi3jeD0GqwP56I7nuPTGbYbIQMCUbG2asCt0hOT20QpMZ_fxcAXiOB6gM-4cCKbIhMjs_xcFS9ldJCR0azabUBufo7tPBQqT2Kp8ttqVNty2a4AClt2-W3DJ1pSnxgC4zEW3cgN1Yqm-wFPOjljmWsdrjp9m5borPSTERknmuOFaPaAYj-aasAv4THb5hwftRloiI6UsjRmR-AZ3BtWTKr9X5QHCxjFpbTGAfo5aFEt1OI1tp7cpyW_m57frFqthet0z2R2HK2FaKj602yR03wRIItC1n91IATiHxHemsllORgXwEjXw6MyRMx59uOZZh7hc_RWX-eq-xJc0u5RJ80qs0kpeE'
--        );



-- INSERT INTO mdb.actors (name, biography, birthdate, origin, profession, avatar)
-- VALUES ('Том Круз',
--         'Том — третий ребёнок в семье, у него есть три сестры: Ли Энн, Мариан и Касс. Ли Энн' ||
--         ' (англ. Lee Anne Devette) сейчас скрывается от популярности брата в Океании. Мать Тома Мэри Ли ' ||
--         'Пфайффер была учителем-дефектологом, а отец Томас Круз Мапотер III — инженером. В поисках ' ||
--         'работы семья часто переезжала, что стало одной из причин трудного детства актёра. Когда Тому ' ||
--         'было 12 лет, его родители развелись. К 14 годам он сменил 15 школ в США и Канаде. Окончательно ' ||
--         'семья Тома осела в Глен Ридж, штат Нью-Джерси, где будущий актёр окончил среднюю школу. В ' ||
--         'детстве у него были кривые, неправильно растущие зубы.\n' ||
--         '\n' ||
--         'Актёрскую карьеру Том начал в 1981 году с небольшой роли в фильме «Бесконечная любовь». ' ||
--         'Свою первую главную роль, которая принесла ему известность, он получил в фильме ' ||
--         '«Рискованный бизнес» (1983).',
--         '3 июля, 1962',
--         'Сиракьюс, Нью-Йорк, США',
--         'актер, продюссер',
--         'https://avatars.mds.yandex.net/get-kinopoisk-image/1946459/2eb2fc4d-a8bd-43b0-83cd-35feacb8ccae/280x420'
--         );

