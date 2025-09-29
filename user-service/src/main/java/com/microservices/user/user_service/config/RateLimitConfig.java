package com.microservices.user.user_service.config;


import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.concurrent.ConcurrentHashMap;

@Configuration
public class RateLimitConfig {

    @Bean
    public ConcurrentHashMap<String, Long> rateLimitCache(){
        return new ConcurrentHashMap<>();
    }
}
