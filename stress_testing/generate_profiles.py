from mimesis import Person, Gender, Address, Text
from dotenv import load_dotenv
import os
import psycopg2
import psycopg2.extras
from collections.abc import Iterator
from uuid import uuid4
import sys, getopt


def mapGender(gender: str) -> str | Gender:
    match gender:
        case 'Муж.':
            return 'male'
        case 'Жен.':
            return 'female'
        case 'male':
            return Gender.MALE
        case 'female':
            return Gender.FEMALE
        case _:
            return None

def newPerson(person: Person, address: Address, text: Text)->tuple[str, str, str, str, str, str, str]:
    gender = mapGender(person.gender())
    return (
        str(uuid4()),
        person.first_name(mapGender(gender)),
        person.last_name(mapGender(gender)),
        person.birthdate().strftime('%Y-%m-%d'),
        gender,
        text.text(),
        address.city()
    )

def generateTuples(n: int, size: int)->Iterator[list[tuple[str, str, str, str, str, str, str]]]:
    locale = 'ru'
    person = Person(locale)
    address = Address(locale)
    text = Text(locale)
    for _ in range(n):
        yield [newPerson(person, address, text) for _ in range(size)]

def printHelp()->None:
    print('generate_profiles.py -n <num batches> -s <batch size>')


def main(argv:list[str])->int:
    load_dotenv('./../.env')
    hostName = 'localhost'
    username = os.environ.get('POSTGRES_USER')
    password = os.environ.get('POSTGRES_PASSWORD')
    databaseName = os.environ.get('POSTGRES_DB')

    psycopg2.extras.register_uuid()

    num_batches = 10
    batch_size = 100000

    try:
        opts, args = getopt.getopt(argv, 'hn:s:',['num-batches=', 'batch-size='])
    except getopt.GetoptError:
        printHelp()
        return 2
    
    for opt, arg in opts:
        if opt == '-h':
            printHelp()
            return 0
        elif opt in ('-n', '--num-batches'):
            num_batches = int(arg)
        elif opt in ('-s', '--batch-size'):
            batch_size = int(arg)
            
    try:
        with psycopg2.connect(dbname=databaseName, user=username, password=password, host=hostName) as conn:
            conn.autocommit = True

            counter = 0
            for tup in generateTuples(num_batches, batch_size):
                with conn.cursor() as curs:
                    args_str = ','.join(curs.mogrify('(%s, %s, %s, %s, %s, %s, %s)', x).decode('utf-8') for x in tup)
                    sql = 'INSERT INTO profiles (id, first_name, last_name, birth_date, gender, biography, city) VALUES ' + args_str
                    curs.execute(sql)
                counter += 1
                print(f'inserted {counter} batches of {num_batches}')
    except Exception as ex:
        print('%s'%ex)

if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
