package com.nchuy099.pastebin_soa.repository;

import com.nchuy099.pastebin_soa.repository.projection.MonthlyStatsProjection;
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
        COALESCE(COUNT(p.id), 0) AS totalPastes,
        COALESCE(SUM(p.views), 0) AS totalViews,
        COALESCE(AVG(p.views), 0) AS avgViewsPerPaste,
        COALESCE(MIN(p.views), 0) AS minViews,
        COALESCE(MAX(p.views), 0) AS maxViews,
        COALESCE(SUM(CASE WHEN p.status = 'ACTIVE' THEN 1 ELSE 0 END), 0) AS activePastes,
        COALESCE(SUM(CASE WHEN p.status = 'EXPIRED' THEN 1 ELSE 0 END), 0) AS expiredPastes
        FROM PasteEntity p
        WHERE YEAR(p.createdAt) = :year 
        AND MONTH(p.createdAt) = :month
        """)

    Optional<MonthlyStatsProjection> getMonthlyStats(@Param("year") int year, @Param("month") int month);
}
