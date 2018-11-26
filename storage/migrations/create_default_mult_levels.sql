CREATE TABLE public.default_mult_levels
(
    number int PRIMARY KEY NOT NULL,
    level json NOT NULL
);
CREATE UNIQUE INDEX default_mult_levels_number_uindex ON public.default_mult_levels (number);