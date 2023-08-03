package com.dhruv.rateLimiter.service;

import com.dhruv.rateLimiter.bean.RateLimiterUser;
import com.dhruv.rateLimiter.repository.RateLimiterRepository;
import org.springframework.stereotype.Service;

import java.time.Instant;
import java.util.logging.Level;
import java.util.logging.Logger;

@Service
public class RateLimiterService {
    private final RateLimiterRepository rateLimiterRepository;
    private final int API_RATE_LIMIT = 7; // Number of API calls per minute
    private final int GLOBAL_USER_RATE_LIMIT = 4; // Total number of API calls per minute allowed for a user
    private final int GLOBAL_API_RATE_LIMIT = 3; // Total number of API calls per minute allowed for any API
    private static final Logger logger = Logger.getLogger(RateLimiterService.class.getName());
    private static RateLimiterUser apiUser;
    private static RateLimiterUser globalUser;
    private static RateLimiterUser globalApiUser;
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
        rateLimiterRepository.save(apiUser);
        rateLimiterRepository.save(globalUser);
        rateLimiterRepository.save(globalApiUser);
        return true;
    }

    private boolean apiSpecificRateLimiter(Long userId, String apiName) {
        apiUser = rateLimiterRepository.findByUserIdAndApiName(userId, apiName);
        Instant now = Instant.now();
        Instant lastRequestTime = (apiUser != null) ? apiUser.getLastRequestTime() : now.minusSeconds(60);
        if (now.isAfter(lastRequestTime.plusSeconds(60)) || 
            now.equals(lastRequestTime.plusSeconds(60))){
            
            // Reset request count if the interval has passed
            if (apiUser == null) {
                apiUser = new RateLimiterUser();
                apiUser.setUserId(userId);
                apiUser.setApiName(apiName);
            }
            apiUser.setLastRequestTime(now);
            apiUser.setRequestCount(1);
        }
        else {
            int requestCount = apiUser.getRequestCount();
            if (requestCount+1 > API_RATE_LIMIT) {
                return false; // Limit exceeded
            }
            apiUser.setRequestCount(requestCount+1);
        }

        logger.log(Level.INFO, "apiSpecifiRateLimiter passed!");
        return true;
    }

    private boolean globalUserRateLimiter(Long userId) {
        globalUser = rateLimiterRepository.findByUserIdAndApiName(userId, "global_user");
        Instant now = Instant.now();
        Instant lastRequestTime = (globalUser != null) ? globalUser.getLastRequestTime() : now.minusSeconds(60);

        if (now.isAfter(lastRequestTime.plusSeconds(60)) || 
            now.equals(lastRequestTime.plusSeconds(60))) {
            
            // Reset request count if the interval has passed
            if (globalUser == null) {
                globalUser = new RateLimiterUser();
                globalUser.setUserId(userId);
                globalUser.setApiName("global_user");
            }
            globalUser.setLastRequestTime(now);
            globalUser.setRequestCount(1);
        } else {
            int requestCount = globalUser.getRequestCount();
            if (requestCount+1 > GLOBAL_USER_RATE_LIMIT) {
                return false; // Limit exceeded
            }
            globalUser.setRequestCount(requestCount+1);
        }
        
        logger.log(Level.INFO, "globalUserRateLimiter passed!");
        return true;
    }

    private boolean globalApiRateLimiter(String apiName) {
        globalApiUser = rateLimiterRepository.findByUserIdAndApiName(-1L, apiName);
        Instant now = Instant.now();
        Instant lastRequestTime = globalApiUser != null ? globalApiUser.getLastRequestTime() : now.minusSeconds(60);

        if (now.isAfter(lastRequestTime.plusSeconds(60))|| 
            now.equals(lastRequestTime.plusSeconds(60))){
            
            // Reset request count if the interval has passed
            if (globalApiUser == null) {
                globalApiUser = new RateLimiterUser();
                globalApiUser.setUserId(-1L);
                globalApiUser.setApiName(apiName);
            }
            globalApiUser.setLastRequestTime(now);
            globalApiUser.setRequestCount(1);
        } else {
            int requestCount = globalApiUser.getRequestCount();
            if (requestCount+1 > GLOBAL_API_RATE_LIMIT) {
                return false; // Limit exceeded
            }
            globalApiUser.setRequestCount(requestCount+1);
        }
        
        logger.log(Level.INFO, "globalApiRateLimiter passed!");
        return true;
    }
}
