CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description VARCHAR(200),
    created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMPTZ
);

CREATE TABLE branches (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE,
    coordinator_name VARCHAR(150),
    contact_number VARCHAR(15) UNIQUE NOT NULL,
    established_on DATE,
    aashram_area NUMERIC,    
    country VARCHAR(100),
    state VARCHAR(100),
    district VARCHAR(100),
    city VARCHAR(100),
    address TEXT,
    pincode VARCHAR(10),
    post_office VARCHAR(100),
    police_station VARCHAR(100),
    open_days VARCHAR(100),         -- e.g., 'Mon-Sun' or 'Mon-Fri'
    daily_start_time TIME,
    daily_end_time TIME,
    parent_branch_id BIGINT NULL,
    created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMPTZ,
    created_by VARCHAR(30),
    updated_by VARCHAR(30)
);

ALTER TABLE branches
    ADD CONSTRAINT fk_branches_parent
    FOREIGN KEY (parent_branch_id) REFERENCES branches(id)
    ON DELETE SET NULL ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_branches_parent ON branches(parent_branch_id);

CREATE TABLE areas (
    id BIGSERIAL PRIMARY KEY,
    branch_id BIGINT NOT NULL REFERENCES branches(id) ON DELETE CASCADE,  -- FK to branches table
    district_id UUID NOT NULL,
    district_coverage FLOAT,
    area_name VARCHAR(100),
    area_coverage FLOAT,
    created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMPTZ,
    created_by VARCHAR(30),
    updated_by VARCHAR(30)
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    contact_number VARCHAR(20),
    password TEXT NOT NULL,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    token TEXT,
    expired_on TIMESTAMPTZ,
    last_login_on TIMESTAMPTZ,
    first_login_on TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMPTZ,
    created_by VARCHAR(30),
    updated_by VARCHAR(30)
);


CREATE TABLE countries (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE states (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    country_id BIGINT NOT NULL REFERENCES countries(id) ON DELETE RESTRICT,
    UNIQUE (name, country_id)
);

-- Create MediaCoverageType table
CREATE TABLE media_coverage_type (
    id SERIAL PRIMARY KEY,
    media_type VARCHAR(50) NOT NULL
);

-- Create PromotionMaterial table
CREATE TABLE promotion_material_type (
    id SERIAL PRIMARY KEY,
    material_type VARCHAR(50) NOT NULL
);

-- -- Need to delete later
-- CREATE TABLE event_details (
--     id SERIAL PRIMARY KEY,
--     event_name VARCHAR(100) NOT NULL
-- );


CREATE TABLE event_types (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);


CREATE TABLE event_categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    event_type_id BIGINT NOT NULL REFERENCES event_types(id)
);

CREATE TABLE event_details (
    id BIGSERIAL PRIMARY KEY,

    event_type_id BIGINT REFERENCES event_types(id),
    event_category_id BIGINT REFERENCES event_categories(id),

    scale VARCHAR(50),
    theme TEXT,

    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    daily_start_time TIME,
    daily_end_time TIME,

    spiritual_orator VARCHAR(200),

    country VARCHAR(100),
    state VARCHAR(100),
    city VARCHAR(100),
    district VARCHAR(100),
    post_office VARCHAR(100),
    pincode VARCHAR(20),
    address TEXT,

    beneficiary_men INT DEFAULT 0,
    beneficiary_women INT DEFAULT 0,
    beneficiary_child INT DEFAULT 0,

    initiation_men INT DEFAULT 0,
    initiation_women INT DEFAULT 0,
    initiation_child INT DEFAULT 0,

    created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMPTZ,
    created_by VARCHAR(100),
    updated_by VARCHAR(100)
);

CREATE TABLE event_media (
    id SERIAL PRIMARY KEY,
    media_coverage_type_id INT NOT NULL REFERENCES media_coverage_type(id),
    event_id INT REFERENCES event_details(id),
    company_name VARCHAR(100) NOT NULL,
    company_email VARCHAR(100),
    company_website VARCHAR(150),
    gender VARCHAR(10),
    prefix VARCHAR(10),
    first_name VARCHAR(50) NOT NULL,
    middle_name VARCHAR(50),
    last_name VARCHAR(50) NOT NULL,
    designation VARCHAR(100),
    contact VARCHAR(20),
    email VARCHAR(100),
    created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMPTZ,
    created_by VARCHAR(30),
    updated_by VARCHAR(30)
);

CREATE TABLE promotion_material_details (
    id SERIAL PRIMARY KEY,
    promotion_material_id INT NOT NULL REFERENCES promotion_material_type(id),
    event_id INT REFERENCES event_details(id),
    quantity INT NOT NULL,
    size VARCHAR(50),
    dimension_height NUMERIC(10,2),
    dimension_width NUMERIC(10,2),
    created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMPTZ,
    created_by VARCHAR(30),
    updated_by VARCHAR(30)
);

CREATE TABLE branch_infrastructure (
    id BIGSERIAL PRIMARY KEY,
    branch_id BIGINT NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    type VARCHAR(100) NOT NULL,
    count INT NOT NULL,
    created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMPTZ,
    created_by VARCHAR(30),
    updated_by VARCHAR(30)
);

CREATE TABLE branch_member (
    id BIGSERIAL PRIMARY KEY,
    member_type VARCHAR(100) NOT NULL,
    name VARCHAR(150) NOT NULL,
    branch_role VARCHAR(100),
    responsibility TEXT,
    age INT,
    date_of_samarpan DATE,
    qualification VARCHAR(150),
    date_of_birth DATE,
    branch_id BIGINT NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMPTZ,
    created_by VARCHAR(30),
    updated_by VARCHAR(30)
);

CREATE TABLE special_guests (
    id BIGSERIAL PRIMARY KEY,
    gender VARCHAR(20),
    prefix VARCHAR(10) NOT NULL,
    first_name VARCHAR(100),
    middle_name VARCHAR(100),
    last_name VARCHAR(100),
    event_id INT REFERENCES event_details(id),
    designation VARCHAR(150),
    organization VARCHAR(200),
    email VARCHAR(150) UNIQUE,
    city VARCHAR(100),
    state VARCHAR(100),
    personal_number VARCHAR(20),
    contact_person VARCHAR(20),
    contact_person_number VARCHAR(20),
    reference_branch_id VARCHAR(100),
    reference_volunteer_id VARCHAR(100),
    reference_person_name VARCHAR(150),
    created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMPTZ,
    created_by VARCHAR(100),
    updated_by VARCHAR(100)
);

CREATE TABLE volunteers (
    id BIGSERIAL PRIMARY KEY,
    branch_id BIGINT NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    volunteer_name VARCHAR(150) NOT NULL,
    number_of_days INTEGER,
    seva_involved VARCHAR(100),
    mention_seva VARCHAR(200),
    event_id INT REFERENCES event_details(id),
    created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMPTZ,
    created_by VARCHAR(100),
    updated_by VARCHAR(100)
);

CREATE TABLE donations (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL,
    branch_id INTEGER NOT NULL,
    donation_type VARCHAR(255),
    amount DOUBLE PRECISION,
    kind_type VARCHAR(255),
    created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMPTZ,
    created_by VARCHAR(100),
    updated_by VARCHAR(100)
);

CREATE TABLE districts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    state_id INTEGER NOT NULL REFERENCES states(id),
    country_id INTEGER NOT NULL REFERENCES countries(id)
);

CREATE TABLE cities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    state_id INTEGER NOT NULL REFERENCES states(id)
);
