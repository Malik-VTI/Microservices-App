package com.microservices.user.user_service.dto;


import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class AuthResponse {

    private String token;
    private String email;
    private String fullName;

    public AuthResponse(String token, String email, String fullName) {
        this.token = token;
        this.email = email;
        this.fullName = fullName;
    }
}
