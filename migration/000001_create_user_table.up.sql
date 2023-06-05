--
-- PostgreSQL database dump
--
-- Dumped from database version 14.8 (Ubuntu 14.8-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.8 (Ubuntu 14.8-0ubuntu0.22.04.1)
SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;
SET default_tablespace = '';
SET default_table_access_method = heap;
--
-- Name: user; Type: TABLE; Schema: public; Owner: user
--
CREATE TABLE public."user" (
                               id character varying(64)  NOT NULL,
                               created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                               updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                               deleted_at timestamp without time zone,
                               first_name character varying(90) NOT NULL,
                               second_name character varying(90) NOT NULL,
                               birthdate character varying(20) NOT NULL,
                               biography text,
                               city character varying(64),
                               password character varying(64)
);
ALTER TABLE public."user" OWNER TO "user";
--
-- Name: user user_pkey; Type: CONSTRAINT; Schema: public; Owner: user
--
ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);
--
-- Name: idx_first_name; Type: INDEX; Schema: public; Owner: user
--
CREATE INDEX idx_first_name ON public."user" USING btree (first_name);
--
-- Name: idx_second_name; Type: INDEX; Schema: public; Owner: user
--
CREATE INDEX idx_second_name ON public."user" USING btree (second_name);
--
-- PostgreSQL database dump complete
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 14.8 (Ubuntu 14.8-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.8 (Ubuntu 14.8-0ubuntu0.22.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: session; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.session (
                                token character varying(64) NOT NULL,
                                user_id character varying(64) NOT NULL,
                                created_at timestamp without time zone,
                                updated_at timestamp without time zone,
                                deleted_at timestamp without time zone,
                                token_till timestamp without time zone
);


ALTER TABLE public.session OWNER TO "user";

--
-- Name: session session_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.session
    ADD CONSTRAINT session_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id);


--
-- PostgreSQL database dump complete
--

