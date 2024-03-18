create table if not exists films(
    id serial,
    title varchar(150) not null check(length(title) > 0),
    description varchar(1000),
    release_date date,
    rating double precision check (rating >= 0 and rating <= 10),
    primary key (id)
);

create table if not exists actors(
    id serial,
    name varchar(100) not null,
    gender varchar(20),
    birth_date date,
    primary key (id)
);

create table if not exists films_x_actors(
    film_id integer,
    actor_id integer,
    foreign key (film_id)
    references films(id),
    foreign key (actor_id)
    references actors(id)
);

create table if not exists users(
    login varchar(50),
    hashed_password varchar(100),
    is_admin boolean,
    primary key (login)
);