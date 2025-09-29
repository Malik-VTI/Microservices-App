package com.microservices.user.user_service.Service;


import com.microservices.user.user_service.model.User;
import com.microservices.user.user_service.repository.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Optional;
import java.util.concurrent.TimeUnit;

@Service
public class UserService {

    @Autowired
    private UserRepository userRepository;

    @Autowired
    private PasswordEncoder passwordEncoder;

    @Autowired
    private RedisTemplate<String, Object> redisTemplate;

    public User createUser(User user) {
        if (userRepository.existsByEmail(user.getEmail())) {
            throw new RuntimeException("User with this email already exists");
        }
        user.setPassword(passwordEncoder.encode(user.getPassword()));
        User savedUser = userRepository.save(user);

        redisTemplate.opsForValue().set("user:" + savedUser.getId(), savedUser,1, TimeUnit.HOURS);
        return savedUser;
    }

    public Optional<User> getUserById(Long id) {
        User cachedUser = (User) redisTemplate.opsForValue().get("user:" + id);
        if (cachedUser != null) {
            return Optional.of(cachedUser);
        }

        Optional<User> user = userRepository.findById(id);
        user.ifPresent(u -> redisTemplate.opsForValue().set("user:" + id, u, 1, TimeUnit.HOURS));
        return user;
    }

    public Optional<User> getUserByEmail(String email) {
        return userRepository.findByEmail(email);
    }

    public List<User> getAllUsers() {
        return userRepository.findAll();
    }

    public User updateUser(Long id, User userDetails) {
        User user = userRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("User not found with id: " + id));

        user.setFirstName(userDetails.getFirstName());
        user.setLastName(userDetails.getLastName());

        User updatedUser = userRepository.save(user);
        redisTemplate.opsForValue().set("user:" + id, updatedUser, 1, TimeUnit.HOURS);
        return updatedUser;
    }

    public void deleteUser(Long id) {
        userRepository.deleteById(id);
        redisTemplate.delete("user:" + id);
    }
}
