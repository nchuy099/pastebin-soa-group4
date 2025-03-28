package com.nchuy099.pastebin_soa.repository;

import com.nchuy099.pastebin_soa.dto.projection.MonthlyStatsProjection;
import com.nchuy099.pastebin_soa.model.PasteEntity;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface PasteRepository extends JpaRepository<PasteEntity, String> {

    @Query("""
            SELECT 
                COUNT(p.id) AS totalPastes,
                SUM(p.views) AS totalViews,
                AVG(p.views) AS avgViewsPerPaste,
                MIN(p.views) AS minViews,
                MAX(p.views) AS maxViews,
                SUM(CASE WHEN p.status = 'ACTIVE' THEN 1 ELSE 0 END) AS activePastes,
                SUM(CASE WHEN p.status = 'EXPIRED' THEN 1 ELSE 0 END) AS expiredPastes
            FROM PasteEntity p
            WHERE YEAR(p.createdAt) = :year 
            AND MONTH(p.createdAt) = :month
            """)
    Optional<MonthlyStatsProjection> getMonthlyStats(@Param("year") int year, @Param("month") int month);
}
