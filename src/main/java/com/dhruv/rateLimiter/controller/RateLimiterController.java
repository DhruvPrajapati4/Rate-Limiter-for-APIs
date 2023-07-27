/*
 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Class.java to edit this template
 */
package com.dhruv.rateLimiter.controller;
import com.dhruv.rateLimiter.service.RateLimiterService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
/**
 *
 * @author dhruv
 */

@RestController
public class RateLimiterController {
    private final RateLimiterService rateLimiterService;

    public RateLimiterController(RateLimiterService rateLimiterService) {
        this.rateLimiterService = rateLimiterService;
    }
    
    @GetMapping("/")
    public String apiDefault() {
        // Implement your GET API logic here
        return "Welcome to rate limiter application! Please enter user id along with request url!";
    }
        
    @GetMapping("/api-alpha/{userId}")
    public ResponseEntity<String> apiAlpha(@PathVariable Long userId) {
        // API-specific rate limiter check for 'API-ALPHA'
        if (!rateLimiterService.allowRequest(userId, "API-ALPHA")) {
            return ResponseEntity.status(HttpStatus.TOO_MANY_REQUESTS)
                    .body("API-ALPHA rate limit exceeded. Please try again later.");
        }

        // API-ALPHA logic here
        return ResponseEntity.ok("API-ALPHA response");
    }

    @PostMapping("/api-beta/{userId}")
    public ResponseEntity<String> apiBeta(@PathVariable Long userId) {
        // API-specific rate limiter check for 'API-BETA'
        if (!rateLimiterService.allowRequest(userId, "API-BETA")) {
            return ResponseEntity.status(HttpStatus.TOO_MANY_REQUESTS)
                    .body("API-BETA rate limit exceeded. Please try again later.");
        }

        // API-BETA logic here
        return ResponseEntity.ok("API-BETA response");
    }

    @PutMapping("/api-gamma/{userId}")
    public ResponseEntity<String> apiGamma(@PathVariable Long userId) {
        // API-specific rate limiter check for 'API-GAMMA'
        if (!rateLimiterService.allowRequest(userId, "API-GAMMA")) {
            return ResponseEntity.status(HttpStatus.TOO_MANY_REQUESTS)
                    .body("API-GAMMA rate limit exceeded. Please try again later.");
        }

        // API-GAMMA logic here
        return ResponseEntity.ok("API-GAMMA response");
    }

    @DeleteMapping("/api-delta/{userId}")
    public ResponseEntity<String> apiDelta(@PathVariable Long userId) {
        // API-specific rate limiter check for 'API-DELTA'
        if (!rateLimiterService.allowRequest(userId, "API-DELTA")) {
            return ResponseEntity.status(HttpStatus.TOO_MANY_REQUESTS)
                    .body("API-DELTA rate limit exceeded. Please try again later.");
        }

        // API-DELTA logic here
        return ResponseEntity.ok("API-DELTA response");
    }
}

