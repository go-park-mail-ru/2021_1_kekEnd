\set ON_ERROR_STOP 1

DROP database mdb;

drop user if exists mdb;

create database mdb;

create user mdb with password 'mdb';

\connect mdb

create schema mdb;

grant usage on schema mdb to mdb;



CREATE TABLE mdb.movie
(
    id              serial PRIMARY KEY,
    title           text not null,
    description     text,
    voiceover       VARCHAR(100)[],
    subtitles       VARCHAR(100)[],
    quality         VARCHAR(20),
    productionYear  integer,
    country         VARCHAR(100)[],
    genre           VARCHAR(100)[],
    slogan          VARCHAR(20),
    director        VARCHAR(20),
    scriptwriter    VARCHAR(20),
    producer        VARCHAR(20),
    operator        VARCHAR(20),
    composer        VARCHAR(20),
    artist          VARCHAR(20),
    montage         VARCHAR(20),
    budget          VARCHAR(20),
    duration        VARCHAR(20),
    actors          VARCHAR (100)[],
    poster          text,
    banner          text,
    trailerPreview  text,

    rating          INTEGER DEFAULT 0
);

GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.movie TO mdb;

COMMENT ON TABLE mdb.movie IS 'Фильмы';

INSERT INTO mdb.movie (title, description, voiceover, subtitles, quality, productionYear, country,
                       genre, slogan, director, scriptwriter, producer, operator, composer, artist,
                       montage, budget, duration, actors, poster, banner, trailerPreview)
VALUES ('Чужой', 'Группа космонавтов высаживается на неизвестной планете и знакомится с ксеноморфом. Шедевр Ридли Скотта',
        '{"Русский", "Английский"}', '{"Русские"}', 'HD', 1979, '{"Великобритания", "США"}',
        '{"Хоррор", "Драма"}', 'slogan', 'Ридли Скотт', '---', '---', '---', '---', '---', '---', '---', '---',
        '{"Сигурни Уивер", "Иэн Холм"}', 'https://avatars.mds.yandex.net/get-kinopoisk-image/1704946/14af6019-b2fe-4e1e-bee5-334d9e472d94/300x450',
        'https://kinopoisk-ru.clstorage.net/E3Y0U4315/f882714f3/JQb5hfuBU5QkfR6H5sDiqWH4CdAOpqGh5AUpIFJuMUGWopdsc4iigraQMDvjNkTICX7J3xVeK_VrWrgjPNqrlmALqJy2hRKhTQhu5t-ku_QU2c0UnnMVQUycuYCmh1BUhv_iLRaJqrvTnVxC6ifhcEGWDVqAjgfY3g9BZlSSplf1Sl3Z190y4Mh9y-faZGYlbfGJsHrf3RU8MFEw5h8klIt7lkES6bq7e0FppppDjcEJSvHiwCpzxHJHtTIsVq7fxcZxZSddQsQIBfsr0ozC_VXNYWxK3xklyBAVWJ5jzNxH9-aFovVSj_OZ3bayy6ERobZp1oAaM4TaDwk61PJSK0nixcVXaaq0BEn7H2bYBgl9kQHsDk8hEfTQiXhmrwhoS1-eyXYBjpPbbZVqgsNdNcV-HXIMdoNUag_lhghqhn-tNjFx0-HG_DxB369K9JpxKfEJdPK3GXEQAEV05heU7Lvv6jV6UYq_S-GZ5nJXNSl1upnW2N5HMJIzbdZMgmIvPb6twSdVPmQkIfOzGugKkSk9DeQiwzX5IFx9DBaHiBCbLwbhsqE-A9tJ3Zoeo5GJnd6Z9piSI8wSrxXOIPLS6wE2bc2nxd7MXN3n_-JManUpLWHUSs-FiaTEtaiyNwSca5OG0YrppiP7bTUq_qcFXdlWEfJUTovY7kPJstSyjsNtwjEZv21GhJQtn98mlFZl6cklrFZn_ZFstCEUhg8U3BfLjsmKldIj1-2tdqrLJdUlojHWcCK34NIPee5UzuaTKT7FWTuZxhAg3VfnulTC6UkFnai2q1G1WASh_LYzJIB_S2Y9GuVOF__hzWYmoxndUSLBGtiaE6SisxFaXGqWi_3Kyc3v4d7EwLHXNxLwmn2xCZkw2i8RpdjwhTS2uxSQ-4veYbYhVk_bEd3echshHSky-SpEessoHovVBnBe2hcB2qkBV-0u-BxFF-_i1EIxfQmh5H6TidlE9MWkplecKB___sk69VZnMxkxQkonGcF5ss3atLbrRPLT0drwkgYvrTq1ed-ZGow0YUcbvhTejVHRcSxK_9m5EEjZNOb3yLjHX3K59lFWJ9-1vWoe4w0ZgeKF4vS261iOv4VKiOqKhzWCZYXrmVaMzEUHe0Z4GvmlTfmAMmeBeaTw3TA61xhAX_9S6RLV9heT4a2aNltZQfFmTfp4pptI2qPtfiTK6p81GjWxd4XeOIBRyxdiyA4teT0hpEbLTV1MQDF4tlsgaM-jxg2qqb4zY-3NVqp7tbkJFpEWwDqPbBp_iXJgfm4HqR41ydcdtjQ4BUOjBuBq7cGB5XjS_40htMQZ1B5rTLQXH_qd5mHiHz9NPRIeZ1EFfTYJ1vh2h0xe08lueA56f3F2tR13AaK8vPUPN5boVrExQQFIhqsJdei81Yw2k0wQM4N-SSoFPm9Xlbne-juVRVn-cQ7QjifU-iu17njqNp-lkuXh36EKvIzlQ3PC0JpdwQGRWH5vCXkYEDX0UqtoDGtPhjW6idbbk-k9xq7DTTnRpkWuEE5_xBID7SZ8TrIfJXaVTceFqmi8_UM31sCGtU2xEcAK59mFRHjJ2JofLERvp5bhsi0mo_v5TaaK7wk1ESadkoTCh1SOl0lOCJK2bzm6xV2ncRpoHCELJ7Z0FhW1zQFkun_dMdhsqZQ6a5Akfy-ubT5p9stf1dVCSk9BHRkSDWZwAmfA8jcNPnAy8mMR0vkxP_kWqFjlE_9CbFIJ0fGpSMLLPSGUhNH0rjsstPsXHi2avYLDq529Ig4rQcEVto3ucAY_7IrbUUo03qYT6XapGUsdGhBULf8jkhyysfEtPbSe40HVfMQJ-FKznDAT61qhFpnmvx8x1ZrypwWF5eL90vBan1wqm1mGqDY-RwHWaUmrjX50rIkfXwr0noFRQSFciq99lbiYRdwmN7BU19cmKf4F2iuvhXWegsOlAXkWQV5gJkv0ZlMBovTCvkeZYkkZaw0KJDxRV3tSzMpdxd3hJCJPSSFsMIlIkpOkrFfjXrXGBaZLyzXJRmbDJQXFNknaDK7zwHYrtX44XnZP-Zo1nZ8dDqTEoVuHBmRCpS2pHXQ218m1-PjNVK6DLMDPU3btoo32AyvZGWr-G421JbJJGhS-h1galz0CvIJqHw3CPT3f9UoUrPHPFz70Tn1pxTUMQnO5MeDwtdSmVyw0wwtezTYlvscnlbFmlps1KclabYpQSjOgPt-RPozuGs-tPh1xJ4n6dDhBe6uWeMrtxYlh4DZLKXWQ-HWQRicENH_PSl163bJTn8V5xkJXkWUd5kFKfALHaLILASrAzr5TPeapFb9tCjwUWRc7opDiLRUh5Xye3919oIR1vCbrbMgDw24hvgm6CxONlfLKby1BXU7FQsCKl4SWx-VqjEb6W2XGcWFzYRZINM2LA6aECl3dvTlE3ru1dbQw3VReN6gkV4NiQTJt5uf_5enmtj-9XYE2UfpEdgM8XgudXlxuAocFPkX5lwXaAFRlY6_C0IKxnalxxEqfLQHU6D1IRjfkBB8LGqVugcofq5lNarY3AYWF6hmOhDoTUO7TPe4k9qZjpQYZdSfdQjyEfU_nZtx2HS0lKaCWpyGl2IxdiEZ3FKSXy8KFSi2iI2s1YeZq2zFNqeppDrh6h0iGP5UqzDbSG0UKXdmjkXooUGFbo-ZI4h1xuXEkOmsJZUycCawqayhcX5umHUY98k_7Bb1mgr85HRV-MZbEVstgdg8VhrjGoiPtjk0Vl_3aIMBBh3_GbIoFabFVsCpPmY2UBBXUFhewXGeXLqlO0aIne-1NXma3zcmZTjlSSBq3bB6HWXr0FhpvmbalUXMZRuRU3XcTEpy--R1dVbCWy_UZVPilgLLvJOzva0YVSjGqH18dKcYKw5m9Bf6x-pTWJ-ieu51WIAJ-X7naeQ27RbqgLOEDL-L81pExhWlIKk9JuWCciQiG59Ckc0uSNfKVJgMnefkiAt_9nUl6nVY8vvM0NnsA',
        'https://kinopoisk-ru.clstorage.net/E3Y0U4315/f882714f3/JQb5hfuBU5QkfR6H5sDiqWH4CdAOpqGh5AUpMCYHMTyD-989NpHas7rlbX6SQwnFJE-UnxlPftgzVp1zDZvnkzQK7JHuiFPIBGAbpt-Zg-gkxTiJQ6LRsM0RXFxjZlz1lzsyjao9Fjs_de3eGi9B0Hle_fZgXhPAdjLlKjwuR8vFimkJb-EiSADxV4PC2EZtYTk9wA7b0WXgzEVEooOQELN_lmkG9UJLk51dnoJbgd39KpHyAE4TaKoHAVpYQvaDGXbRScuVPuRIBQP7BtDOXe0VgTCOq_WlmPyJTEonSCiPl0pZBjFyv7PZNeb-s8kRcVpVQhymJ-BysxFy0G4mK-Xu5Qn71X5MALkLA6IsOtHJoSFAukOFqdT8vfge47wA0_fKHZYJRoNzRRUeMnd9HdkiGZboejvIhvvNQtiSAk8tmkGFl2Wu_ABVt986-G41HVWVxDIvzZHc2A3UMn_EVLP_UuFuMSZno4kZnorr-e1pVuWauH4XWAr3QbJM-h7_6R7lUT-hygQ83eeb0hw-pZXJJTCG11WJvPQRhC7vJFR_-8IxGh2uV_-xYW4a9wXNqboxUkjKq8iisz3euF6-F40exRFzxTaEoKX_96Ks-pkdCSmwFjf98UDgXXyKj8BUb5fyWRat3tNnecXKqnuJPRU6AWoAUovohluNPrzKalMFlsHdr20ieEx5gwPi8NL16VX9qL6vtTmwmH24UrPQlIvTikW-1cYbF-15_pJryUFNMhlqfCaLxAbD0WrQ6uKv8bbl-cNRGkQAyd-DnpiCsRWhvSxKLyGNQFBFJJInXDRHa46lcnniL6dtkc6uW5Upzdrt-gy6v-wKo8HmuNbq23FGKVF79V40vKFrizrommXhrSn4Mjf1bSzQlYw2f8jMS28WyfY58q9TSVnOJmuFrQ1isVbIoufI-rN5sgDuKqNhfhnNmy3SiIRlN6cOpAaZ8c3lQD7HybHYEE18EqeEAEtXwm1KoY4zVwnJ3srjPUl5QhnaHKLPIPJf5Yo81vbz4UrpPVcNvmTcYesnwng-NRHRnchK872Z_EC5INI7OCyThwpZJvHuZ-sVWZ5qt62R2c5pFriij8xe083e-MIuC7EC0X1XDaIYsDV7X7r0lq2pAWH8Sr-9YdgA2di-_0zYDw-mIb6pLtNTEV1CSmdVCXnuOfI8Ar-ACsM99kCWdns1ysnxR32yTKxdT_OalI6tMVFVYFY3CS3MzLX8DuuYBH_XglUSZQo74_0VzsZffZklet1KQEqbcAaj8Wpgeo6DRRYlSdtploxwOUO3LhAWMTVRLcDOXwWVmEQBmCaPWLzDE17BJqV2w2fVuWb2M6FBmUZNBogWtyymU7XefJ4y92WO5XGXYbbI3Hlfr14Ebuld0flg0kuNEWgIlQguswRMA_dulXIhIp8fGeFODjMFZQXCmcrsysdEftd5OiBactOt9j1Zb8EubCQF36-6SI49uYEFyHLjjSF4RNFcFn_ovENnfm22IS5vs_mZKjYXGT3JOuVaYCJzgAJTYW7Ygg5b9cKdma-ZPoQMXRerHswOvV3xqdBWQ1kRYESVSAZjADDz5-YZPvHSM9sFteKCU1E5ISoxUsTSC-gSIwFK9MYCm3UaoQ0jYa4YmPl_38LIfqGRobmwovNZsbwMhSiy86DIj_dCqab1Zq_PZflC9u8xKakSvd6AAmNMPrvlilSOKpNBilX544E6ZDi9D6dijHKJ-Z3VKCr_mfV4FF3cqre8rLNfbtESFXbjJx2Z1qZToa2Rov16VHZruHbThc4wjvaf5Qrd-efZFhzU4XvjjtgCcV3N_VzO8yH5sPiBDNpXBIxvy5KNOmmCC2fFlSou4yVFbeZx9nASFwzauz0yvMqyb7F64Xm7eaa8lOm3f2ZAVpn9Da28XpdFARQY_ZQyezQsA9d6mXZVws87ibFeqs9BgVGa-R7sLoO8bhs5QthqNvNFxm3px60O8FyxkyOSwFYBSS39fN7jFZHMUNnMCi_ouJ8XAjGWYXYbk0Ul6g7buQFl4mUm7FLj2N6n4abY6jJPZc7phU8VOuAkBU_vDgheYbFReYjO55VpPFwlmKKnEFDr61IlDuHij1sBOdYeU9WZ1co9QmQCqzgyd80-AEKCr-HOKZ1fYaKMmI0za9IUDpXpWdnIJqMlAWzItaAyq8gUh8MqUaqRZpdTebneylMhlY3iHdbMSm80ft_BVoD6HkMJ6rnZq9VaqNAhD1u-ZN41FXmVMFoTRZXcfAkIvi9YuMuXxiWSASLnW7n9Prp7ISlJ9o2aNEb7jC4XYYJMXlKXtcZ59eMhkiQEsRsXnsBCpc3N8ai-4w25xBCZPFYHmGhjE1qNBvUq1ye50V52E91VRdLxXuBOowBm-1UKdOJ21x1CcUlrcX4AyFVbWxaESv3tFYVksv95mVCMoThC7-ig_89izWKdIsOTETkmqtcxAQXekdKEEk_sDodBdiRyagtl1snNl-XGyAQtb4s-fJadFSEdgNYzMfn4ZA1cFmcE4OuH4llGBVajS_ElPqqbEUmNpnWOaD63uHIjzXYszrIPuZ69Ddv1qnjcjd_zpthyPS19kTAOqw0p4EhF-BqTqFBn34aFfgnyry-R5T7qa7HBTX5VqsRWi3jeD0GqwP56I7nuPTGbYbIQMCUbG2asCt0hOT20QpMZ_fxcAXiOB6gM-4cCKbIhMjs_xcFS9ldJCR0azabUBufo7tPBQqT2Kp8ttqVNty2a4AClt2-W3DJ1pSnxgC4zEW3cgN1Yqm-wFPOjljmWsdrjp9m5borPSTERknmuOFaPaAYj-aasAv4THb5hwftRloiI6UsjRmR-AZ3BtWTKr9X5QHCxjFpbTGAfo5aFEt1OI1tp7cpyW_m57frFqthet0z2R2HK2FaKj602yR03wRIItC1n91IATiHxHemsllORgXwEjXw6MyRMx59uOZZh7hc_RWX-eq-xJc0u5RJ80qs0kpeE'
        );


CREATE TABLE mdb.users
(
    login               VARCHAR(100) PRIMARY KEY,
    password            VARCHAR(256) NOT NULL,
    img_src             text DEFAULT 'http://89.208.198.186:8080/avatars/default.jpeg',

    firstname           VARCHAR(100),
    lastname            VARCHAR(100),
    sex                 INTEGER CONSTRAINT sex_t CHECK (sex = 1 OR sex = 0),
    email               VARCHAR(100) NOT NULL UNIQUE,
    registration_date   timestamp NOT NULL DEFAULT NOW(),

    description         VARCHAR(600),
    movies_watched      INTEGER DEFAULT 0,
    reviews_count       INTEGER DEFAULT 0,
    friends_count       INTEGER DEFAULT 0,
    user_rating         INTEGER DEFAULT 0
);

GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.users TO mdb;

COMMENT ON TABLE mdb.users IS 'Пользователи';




-- TO DO Сделать тригер на пересчет рейтинга в поле rating таблицы mdb.movie
CREATE TABLE mdb.movie_rating
(
    user_login VARCHAR(100) REFERENCES mdb.users (login) ON DELETE CASCADE,
    movie_id INTEGER REFERENCES mdb.movie (id) ON DELETE CASCADE,
    rating INTEGER CONSTRAINT from_one_to_ten_rating CHECK (rating > 1 AND rating <= 10) NOT NULL,
    PRIMARY KEY (user_login, movie_id)
);

GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.movie_rating TO mdb;

COMMENT ON TABLE mdb.movie_rating IS 'Рейтинг фильмов';

CREATE OR REPLACE FUNCTION rating_recalc() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE mdb.movie SET rating=rating + NEW.rating WHERE id=NEW.movie_id;
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        UPDATE mdb.movie SET rating=rating - OLD.rating + NEW.rating WHERE id=NEW.movie_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE mdb.movie SET rating=rating - OLD.rating WHERE id=OLD.movie_id;
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

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA mdb TO mdb;
