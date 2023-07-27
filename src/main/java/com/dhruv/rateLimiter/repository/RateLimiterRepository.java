/*
 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Interface.java to edit this template
 */
package com.dhruv.rateLimiter.repository;
import com.dhruv.rateLimiter.bean.RateLimiterUser;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
/**
 *
 * @author dhruv
 */


@Repository
public interface RateLimiterRepository extends JpaRepository<RateLimiterUser, Long> {
    RateLimiterUser findByUserIdAndApiName(Long userId, String apiName);
}

