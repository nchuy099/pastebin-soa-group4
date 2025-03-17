package com.nchuy099.pastebin_soa.dto;

import lombok.AllArgsConstructor;
import lombok.Getter;


@Getter
@AllArgsConstructor
public class MonthlyStatsDTO {
        String month;
        int totalPastes;
        int totalViews;
        double avgViewsPerPaste;
        int minViews;
        int maxViews;
        int activePastes;
        int expiredPastes;
}