export type DBConfig = {
    host: string;
    port: number;
    user: string;
    password: string;
    database: string;
    multipleStatements: boolean;
};

export type JWTConfig = {
    secret: string;
};

export const LocalEnv = {
    db: {
        host: process.env.DB_ADDRESS!,
        port: parseInt(process.env.DB_PORT!, 10),
        user: process.env.DB_USER!,
        password: process.env.DB_PASSWORD!,
        database: process.env.DB_NAME!,
        multipleStatements: true,
        supportBigNumbers: true,
        bigNumberStrings: true,
    },
    jwt: {
        secret: process.env.JWT_SECRET!,
    }
  };