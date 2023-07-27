package com.dhruv.rateLimiter.service;

import com.dhruv.rateLimiter.bean.RateLimiterUser;
import com.dhruv.rateLimiter.repository.RateLimiterRepository;
import org.springframework.stereotype.Service;

import java.time.Instant;

@Service
public class RateLimiterService {
    private final RateLimiterRepository rateLimiterRepository;
    private final int API_RATE_LIMIT = 10; // Number of API calls per minute
    private final int GLOBAL_USER_RATE_LIMIT = 6; // Total number of API calls per minute allowed for a user
    private final int GLOBAL_API_RATE_LIMIT = 4; // Total number of API calls per minute allowed for any API

    public RateLimiterService(RateLimiterRepository rateLimiterRepository) {
        this.rateLimiterRepository = rateLimiterRepository;
    }

    public boolean allowRequest(Long userId, String apiName) {
        // Perform API-specific rate limiter check
        if (!apiSpecificRateLimiter(userId, apiName)) {
            return false;
        }

        // Perform Global User Rate Limiter check
        if (!globalUserRateLimiter(userId)) {
            return false;
        }

        // Perform Global API Rate Limiter check
        if (!globalApiRateLimiter(apiName)) {
            return false;
        }

        return true;
    }

    private boolean apiSpecificRateLimiter(Long userId, String apiName) {
        RateLimiterUser userEntry = rateLimiterRepository.findByUserIdAndApiName(userId, apiName);
        Instant now = Instant.now();
        Instant lastRequestTime = (userEntry != null) ? userEntry.getLastRequestTime() : now.minusSeconds(60);
        if (now.isAfter(lastRequestTime.plusSeconds(60)) || 
            now.equals(lastRequestTime.plusSeconds(60))){
            
            // Reset request count if the interval has passed
            if (userEntry == null) {
                userEntry = new RateLimiterUser();
                userEntry.setUserId(userId);
                userEntry.setApiName(apiName);
            }
            userEntry.setLastRequestTime(now);
            userEntry.setRequestCount(1);
        }
        else {
            int requestCount = userEntry.getRequestCount() + 1;
            if (requestCount > API_RATE_LIMIT) {
                return false; // Limit exceeded
            }
            userEntry.setRequestCount(requestCount);
        }

        rateLimiterRepository.save(userEntry);
        return true;
    }

    private boolean globalUserRateLimiter(Long userId) {
        RateLimiterUser userEntry = rateLimiterRepository.findByUserIdAndApiName(userId, "global_user");
        Instant now = Instant.now();
        Instant lastRequestTime = (userEntry != null) ? userEntry.getLastRequestTime() : now.minusSeconds(60);

        if (now.isAfter(lastRequestTime.plusSeconds(60)) || 
            now.equals(lastRequestTime.plusSeconds(60))) {
            
            // Reset request count if the interval has passed
            if (userEntry == null) {
                userEntry = new RateLimiterUser();
                userEntry.setUserId(userId);
                userEntry.setApiName("global_user");
            }
            userEntry.setLastRequestTime(now);
            userEntry.setRequestCount(1);
        } else {
            int requestCount = userEntry.getRequestCount() + 1;
            if (requestCount > GLOBAL_USER_RATE_LIMIT) {
                return false; // Limit exceeded
            }
            userEntry.setRequestCount(requestCount);
        }

        rateLimiterRepository.save(userEntry);
        return true;
    }

    private boolean globalApiRateLimiter(String apiName) {
        RateLimiterUser userEntry = rateLimiterRepository.findByUserIdAndApiName(-1L, apiName);
        Instant now = Instant.now();
        Instant lastRequestTime = userEntry != null ? userEntry.getLastRequestTime() : now.minusSeconds(60);

        if (now.isAfter(lastRequestTime.plusSeconds(60))|| 
            now.equals(lastRequestTime.plusSeconds(60))){
            
            // Reset request count if the interval has passed
            if (userEntry == null) {
                userEntry = new RateLimiterUser();
                userEntry.setUserId(-1L);
                userEntry.setApiName(apiName);
            }
            userEntry.setLastRequestTime(now);
            userEntry.setRequestCount(1);
        } else {
            int requestCount = userEntry.getRequestCount() + 1;
            if (requestCount > GLOBAL_API_RATE_LIMIT) {
                return false; // Limit exceeded
            }
            userEntry.setRequestCount(requestCount);
        }

        rateLimiterRepository.save(userEntry);
        return true;
    }
}
