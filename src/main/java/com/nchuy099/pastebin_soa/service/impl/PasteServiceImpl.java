package com.nchuy099.pastebin_soa.service.impl;

import com.nchuy099.pastebin_soa.common.Status;
import com.nchuy099.pastebin_soa.common.Visibility;
import com.nchuy099.pastebin_soa.dto.MonthlyStatsDTO;
import com.nchuy099.pastebin_soa.model.PasteEntity;
import com.nchuy099.pastebin_soa.repository.PasteRepository;
import com.nchuy099.pastebin_soa.service.PasteService;
import lombok.Getter;
import lombok.RequiredArgsConstructor;
import lombok.Setter;
import org.springframework.stereotype.Service;

import java.text.SimpleDateFormat;
import java.time.Instant;
import java.util.Date;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class PasteServiceImpl implements PasteService {

    private final PasteRepository pasteRepository;

    @Override
    public PasteEntity createPaste(PasteEntity pasteData) throws Exception {
        try {
            String id;
            do {
                id = UUID.randomUUID().toString().substring(0, 8);
            } while (pasteRepository.existsById(id)); // Kiểm tra trùng lặp

            pasteData.setId(id);

            if (pasteData.getTitle() == null)
                pasteData.setTitle("Untitled");
            if (pasteData.getLanguage() == null)
                pasteData.setLanguage("text");
            if (pasteData.getVisibility() == null)
                pasteData.setVisibility(Visibility.PUBLIC);
            if (pasteData.getStatus() == null)
                pasteData.setStatus(Status.ACTIVE);

            return pasteRepository.save(pasteData);
        } catch (Exception e) {
            System.err.println("Create paste error: " + e.getMessage());
            throw new Exception("Failed to create paste");
        }
    }

    @Override
    public PasteEntity getPasteById(String id) throws Exception {
        try {
            Optional<PasteEntity> pasteOptional = pasteRepository.findById(id);
            if (pasteOptional.isEmpty()) {
                throw new Exception("Paste not found");
            }
            PasteEntity paste = pasteOptional.get();

            if (paste.getExpiresAt() != null && paste.getExpiresAt().before(Date.from(Instant.now()))) {
                paste.setStatus(Status.EXPIRED);
                pasteRepository.save(paste);
                throw new Exception("This paste has expired and is no longer accessible");
            }

            if (paste.getStatus() == Status.ACTIVE) {
                pasteRepository.incrementViews(id);
            }

            return paste;
        } catch (Exception e) {
            System.err.println("Get paste by id error: " + e.getMessage());
            throw e;
        }
    }

    @Override
    public List<PasteEntity> getPublicPastes() throws Exception {
        try {
            return pasteRepository.findPublicPastes();
        } catch (Exception e) {
            System.err.println("Get public pastes error: " + e.getMessage());
            throw new Exception("Failed to fetch public pastes");
        }
    }

    @Override
    public MonthlyStatsDTO getMonthlyStats(Date month) throws Exception {
        try {
            Object[][] res = pasteRepository.getMonthlyStats(month);
            Object[] stats = res[0];

            String strMonth = new SimpleDateFormat("yyyy-MM").format(month);

            // Safely convert each element to the appropriate type
            int totalPastes = stats[0] != null ? ((Number) stats[0]).intValue() : 0;
            int totalViews = stats[1] != null ? ((Number) stats[1]).intValue() : 0;
            double avgViewsPerPaste = stats[2] != null ? ((Number) stats[2]).doubleValue() : 0.0;
            int minViews = stats[3] != null ? ((Number) stats[3]).intValue() : 0;
            int maxViews = stats[4] != null ? ((Number) stats[4]).intValue() : 0;
            int activePastes = stats[5] != null ? ((Number) stats[5]).intValue() : 0;
            int expiredPastes = stats[6] != null ? ((Number) stats[6]).intValue() : 0;

            return new MonthlyStatsDTO(
                    strMonth,
                    totalPastes,
                    totalViews,
                    avgViewsPerPaste,
                    minViews,
                    maxViews,
                    activePastes,
                    expiredPastes);
        } catch (Exception e) {
            System.err.println("Get monthly stats error: " + e.getMessage());
            e.printStackTrace(); // Add stack trace for better debugging
            throw new Exception("Failed to fetch monthly statistics");
        }
    }
}