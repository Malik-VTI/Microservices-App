package com.microservices.user.user_service.controller;


import com.microservices.user.user_service.Service.UserService;
import com.microservices.user.user_service.dto.AuthRequest;
import com.microservices.user.user_service.dto.AuthResponse;
import com.microservices.user.user_service.dto.RegisterRequest;
import com.microservices.user.user_service.model.User;
import com.microservices.user.user_service.util.JwtUtil;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/auth")
@CrossOrigin(origins = "*")
public class AuthenticationController {

    @Autowired
    private AuthenticationManager authenticationManager;

    @Autowired
    private UserService userService;

    @Autowired
    private JwtUtil jwtUtil;

    @Autowired
    private PasswordEncoder passwordEncoder;

    @PostMapping("/login")
    public ResponseEntity<?> login(@Valid @RequestBody AuthRequest request) {
        try {
            Authentication authentication = authenticationManager.authenticate(
                    new UsernamePasswordAuthenticationToken(request.getEmail(), request.getPassword())
            );

            User user = userService.getUserByEmail(request.getEmail())
                    .orElseThrow(() -> new BadCredentialsException("Invalid credentials"));

            String token = jwtUtil.generateToken(user.getEmail(), "ROLE_USER");

            return ResponseEntity.ok(new AuthResponse(token, user.getEmail(), user.getFirstName() + " " + user.getLastName()));
        } catch (BadCredentialsException e) {
            return ResponseEntity.badRequest().body("Invalid credentials");
        }
    }

    @PostMapping("/register")
    public ResponseEntity<?> register(@Valid @RequestBody RegisterRequest request) {
        try {
            User user = new User(
                    request.getEmail(),
                    request.getPassword(),
                    request.getFirstName(),
                    request.getLastName()
            );

            User createdUser = userService.createUser(user);
            String token = jwtUtil.generateToken(createdUser.getEmail(), "ROLE_USER");

            return ResponseEntity.ok(new AuthResponse(token, createdUser.getEmail(),
                    createdUser.getFirstName() + " " + createdUser.getLastName()));
        } catch (RuntimeException e) {
            return ResponseEntity.badRequest().body(e.getMessage());
        }
    }

    @PostMapping("/validate")
    public ResponseEntity<?> validateToken(@RequestHeader("Authorization") String token) {
        try {
            String jwtToken = token.replace("Bearer ", "");
            if (jwtUtil.validateToken(jwtToken)) {
                String username = jwtUtil.extractUsername(jwtToken);
                return ResponseEntity.ok().body("Token is valid for user: " + username);
            } else {
                return ResponseEntity.badRequest().body("Invalid token");
            }
        } catch (Exception e) {
            return ResponseEntity.badRequest().body("Invalid token");
        }
    }
}
