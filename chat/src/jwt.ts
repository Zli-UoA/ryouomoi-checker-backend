import * as jwt from 'jsonwebtoken';
import { JWTConfig } from "./config";

export class UserJWTService {
    secret: string;

    constructor(jwtConfig: JWTConfig) {
        this.secret = jwtConfig.secret;
    }

    getUserIdFromJWT(token: string): string {
        const decoded = jwt.verify(token, this.secret) as jwt.JwtPayload;
        return decoded.userID;
    }
}
