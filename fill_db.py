import psycopg2
import requests
import time

with open('keys.txt') as keys:
    KINO_API_KEY = keys.readline().strip()
    YOUTUBE_API_KEY = keys.readline().strip()

MOVIE_API_PATH = 'https://kinopoiskapiunofficial.tech/api/v2.1/films'
MOVIE_API_APPEND = '?append_to_response=BUDGET&append_to_response=RATING'
STAFF_API_PATH = 'https://kinopoiskapiunofficial.tech/api/v1/staff'
YOUTUBE_TRAILER_PATH = f'https://www.googleapis.com/youtube/v3/search' \
                       f'?key={YOUTUBE_API_KEY}&part=snippet&type=video&maxResults=1&q='
YOUTUBE_EMBED = 'https://www.youtube.com/embed/'

START_MOVIE_INDEX = 300
END_MOVIE_INDEX = 320


def format_array(arr):
    return '{{{0}}}'.format(', '.join([f'"{x}"' for x in arr]))


def filter_by_profession(people, profession):
    return [person['nameRu'] for person in people if person['professionKey'] == profession]


def get_staff_info(movie_id):
    staff_array = requests.get(f'{STAFF_API_PATH}?filmId={movie_id}',
                               headers={'X-API-KEY': KINO_API_KEY}).json()
    return [
        filter_by_profession(staff_array, 'DIRECTOR')[0],
        filter_by_profession(staff_array, 'WRITER')[0],
        filter_by_profession(staff_array, 'PRODUCER')[0],
        filter_by_profession(staff_array, 'OPERATOR')[0],
        filter_by_profession(staff_array, 'COMPOSER')[0],
        filter_by_profession(staff_array, 'DESIGN')[0],
        filter_by_profession(staff_array, 'EDITOR')[0],
        [(person['staffId'], person['nameRu']) for person in staff_array if person['professionKey'] == 'ACTOR']
    ]


def get_frames_info(movie_id):
    frames_array = requests.get(f'{MOVIE_API_PATH}/{movie_id}/frames',
                                headers={'X-API-KEY': KINO_API_KEY}).json()
    frame = frames_array['frames'][0]['image']
    return frame


def get_trailer_info(title):
    trailer_response = requests.get(f'{YOUTUBE_TRAILER_PATH}{title}+трейлер').json()
    trailer_link = YOUTUBE_EMBED + trailer_response['items'][0]['id']['videoId']
    return trailer_link


def get_actor_info(actor_id):
    response = requests.get(f'{STAFF_API_PATH}/{actor_id}',
                            headers={'X-API-KEY': KINO_API_KEY}).json()
    return [
        response['nameRu'],
        "\n\n".join(response['facts']),
        response['birthday'],
        response['birthplace'],
        response['profession'],
        response['posterUrl']
    ]


def get_movie_info(movie_id):
    response = requests.get(f'{MOVIE_API_PATH}/{movie_id}{MOVIE_API_APPEND}',
                            headers={'X-API-KEY': KINO_API_KEY}).json()

    data = response['data']
    rating = response['rating']
    budget = response['budget']

    staff = get_staff_info(movie_id)
    frame = get_frames_info(movie_id)
    trailer_thumbnail = get_trailer_info(data['nameRu'])

    actors = staff[-1][:5]

    return [
        data['nameRu'],
        data['description'],
        data['year'],
        format_array([item['country'] for item in data['countries']]),
        data['slogan'],
        *staff[:-1],
        budget['budget'],
        data['filmLength'],
        data['posterUrlPreview'],
        frame,
        trailer_thumbnail,
        rating['rating'],
        rating['ratingVoteCount']
    ], set([item['genre'] for item in data['genres']]), [actor[0] for actor in actors]


def fill_db(conn, cursor):
    counter = 0
    genres_current_id = 1

    for index in range(START_MOVIE_INDEX, END_MOVIE_INDEX):
        try:
            info, genres, actors_ids = get_movie_info(index)
            info = [item if item is not None else 'нет данных' for item in info]

            cursor.execute(
                'INSERT INTO movie (title, description, productionyear, country, slogan, '
                'director, scriptwriter, producer, operator, composer, artist, montage, '
                'budget, duration, poster, banner, trailerpreview, rating, rating_count) '
                'VALUES(%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)',
                info)

            for genre in genres:
                cursor.execute('SELECT id FROM genres WHERE name = %s', (genre,))
                genre_id = cursor.fetchone()
                if genre_id is None:
                    # этот жанр еще не был добавлен
                    genre_id = genres_current_id
                    cursor.execute('INSERT INTO genres (id, name) VALUES(%s, %s)', (genre_id, genre))
                    genres_current_id += 1

                cursor.execute(
                    'INSERT INTO movie_genres (movie_id, genre_id) VALUES(%s, %s)', (counter + 1, genre_id)
                )

            for actor_id in actors_ids:
                cursor.execute('SELECT id FROM actors WHERE id = %s', (actor_id,))
                if cursor.fetchone() is None:
                    # этот актер еще не был добавлен
                    actor_info = get_actor_info(actor_id)
                    cursor.execute(
                        'INSERT INTO actors (id, name, biography, birthdate, origin, profession, avatar) '
                        'VALUES(%s, %s, %s, %s, %s, %s, %s)', (actor_id, *actor_info)
                    )

                cursor.execute(
                    'INSERT INTO movie_actors (movie_id, actor_id) VALUES(%s, %s)', (counter + 1, actor_id)
                )
        except Exception as e:
            print('id:', index, 'error:', e)
            print(e)
            conn.rollback()
            continue

        counter += 1
        conn.commit()
        time.sleep(0.2)

    cursor.execute(
        'INSERT INTO meta (version, movies_count, users_count)'
        'VALUES(%s, %s, %s)', (1, counter, 0))
    conn.commit()


def main():
    conn = psycopg2.connect(
        host='localhost',
        user='mdb',
        password='mdb',
        database='mdb'
    )
    cursor = conn.cursor()
    fill_db(conn, cursor)
    cursor.close()
    conn.close()


if __name__ == '__main__':
    main()
