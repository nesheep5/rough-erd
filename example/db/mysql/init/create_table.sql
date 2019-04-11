use test;

CREATE TABLE employees (
    id          INT             NOT NULL,
    birth_date  DATE            NOT NULL,
    first_name  VARCHAR(14)     NOT NULL,
    last_name   VARCHAR(16)     NOT NULL,
    gender      ENUM ('M','F')  NOT NULL,
    hire_date   DATE            NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE departments (
    id          INT             NOT NULL,
    dept_name   VARCHAR(40)     NOT NULL,
    PRIMARY KEY (id),
    UNIQUE  KEY (dept_name)
);

CREATE TABLE dept_manager (
   id            INT             NOT NULL,
   employee_id   INT             NOT NULL,
   department_id INT             NOT NULL,
   from_date     DATE            NOT NULL,
   to_date       DATE            NOT NULL,
   FOREIGN KEY (employee_id)  REFERENCES employees (id)    ON DELETE CASCADE,
   FOREIGN KEY (department_id) REFERENCES departments (id) ON DELETE CASCADE,
   PRIMARY KEY (id)
);

CREATE TABLE dept_emp (
    id            INT             NOT NULL,
    employee_id   INT             NOT NULL,
    department_id INT             NOT NULL,
    from_date     DATE            NOT NULL,
    to_date       DATE            NOT NULL,
    FOREIGN KEY (employee_id)  REFERENCES employees (id)    ON DELETE CASCADE,
    FOREIGN KEY (department_id) REFERENCES departments (id) ON DELETE CASCADE,
    PRIMARY KEY (id)
);

CREATE TABLE titles (
    id            INT             NOT NULL,
    employee_id   INT             NOT NULL,
    title         VARCHAR(50)     NOT NULL,
    from_date     DATE            NOT NULL,
    to_date       DATE,
    FOREIGN KEY (employee_id)  REFERENCES employees (id)    ON DELETE CASCADE,
    PRIMARY KEY (id)
);

CREATE TABLE salaries (
    id            INT             NOT NULL,
    employee_id   INT             NOT NULL,
    salary        INT             NOT NULL,
    from_date     DATE            NOT NULL,
    to_date       DATE            NOT NULL,
    FOREIGN KEY (employee_id)  REFERENCES employees (id)    ON DELETE CASCADE,
    PRIMARY KEY (id)
);
