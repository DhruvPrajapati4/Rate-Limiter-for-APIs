package com.dhruv.rateLimiter.bean;

/*
 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Class.java to edit this template
 */

/**
 *
 * @author dhruv
 */
import javax.persistence.*;
import java.time.Instant;

@Entity
@Table(name = "user_details")
public class RateLimiterUser{
    
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "user_id")
    private Long userId;

    @Column(name = "api_name")
    private String apiName;

    @Column(name = "last_request_time")
    private Instant lastRequestTime;

    @Column(name = "request_count")
    private int requestCount;
    
    // Default constructor (required by JPA)
    public RateLimiterUser() {
    }

    // Constructor for API-specific Rate Limiter
    public RateLimiterUser(Long userId, String apiName, Instant lastRequestTime, int requestCount) {
        this.userId = userId;
        this.apiName = apiName;
        this.lastRequestTime = lastRequestTime;
        this.requestCount = requestCount;
    }

    // Constructor for Global User Rate Limiter
    public RateLimiterUser(Long userId, Instant lastRequestTime, int requestCount) {
        this.userId = userId;
        this.apiName = "global_user";
        this.lastRequestTime = lastRequestTime;
        this.requestCount = requestCount;
    }

    // Constructor for Global API Rate Limiter
    public RateLimiterUser(String apiName, Instant lastRequestTime, int requestCount) {
        this.userId = -1L; // Use a fixed value to represent the Global API Rate Limiter
        this.apiName = apiName;
        this.lastRequestTime = lastRequestTime;
        this.requestCount = requestCount;
    }

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public Long getUserId() {
        return userId;
    }

    public void setUserId(Long userId) {
        this.userId = userId;
    }

    public String getApiName() {
        return apiName;
    }

    
    public void setApiName(String apiName) {
        this.apiName = apiName;
    }

    
    public Instant getLastRequestTime() {
        return lastRequestTime;
    }

    
    public void setLastRequestTime(Instant lastRequestTime) {
        this.lastRequestTime = lastRequestTime;
    }

    
    public int getRequestCount() {
        return requestCount;
    }

    
    public void setRequestCount(int requestCount) {
        this.requestCount = requestCount;
    }    
}
