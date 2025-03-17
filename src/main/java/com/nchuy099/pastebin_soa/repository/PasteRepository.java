package com.nchuy099.pastebin_soa.repository;

import com.nchuy099.pastebin_soa.common.Visibility;
import com.nchuy099.pastebin_soa.model.PasteEntity;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;
import org.springframework.transaction.annotation.Transactional;

import java.util.Date;
import java.util.List;

@Repository
public interface PasteRepository extends JpaRepository<PasteEntity, String> {


    // Lấy danh sách paste công khai, chưa hết hạn
    @Query("SELECT p FROM PasteEntity p WHERE p.visibility = 'PUBLIC' AND p.status = 'ACTIVE' ORDER BY p.createdAt DESC")
    List<PasteEntity> findPublicPastes();

    // Tăng lượt xem
    @Modifying
    @Transactional
    @Query("UPDATE PasteEntity p SET p.views = p.views + 1 WHERE p.id = :id")
    void incrementViews(String id);

    // Thống kê theo tháng
    @Query("SELECT " +
            "COUNT(p.id) AS totalPastes, " +
            "SUM(p.views) AS totalViews, " +
            "AVG(p.views) AS avgViewsPerPaste, " +
            "MIN(p.views) AS minViews, " +
            "MAX(p.views) AS maxViews, " +
            "SUM(CASE WHEN p.status = 'ACTIVE' THEN 1 ELSE 0 END) AS activePastes, " +
            "SUM(CASE WHEN p.status = 'EXPIRED' THEN 1 ELSE 0 END) AS expiredPastes " +
            "FROM PasteEntity p " +
            "WHERE YEAR(p.createdAt) = YEAR(:month) " +
            "AND MONTH(p.createdAt) = MONTH(:month)")
    Object[][] getMonthlyStats(Date month);
}
