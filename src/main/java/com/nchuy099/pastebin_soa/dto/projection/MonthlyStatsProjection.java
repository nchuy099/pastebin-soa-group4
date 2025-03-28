package com.nchuy099.pastebin_soa.dto.projection;

public interface MonthlyStatsProjection {
        Long getTotalPastes();

        Long getTotalViews();

        Double getAvgViewsPerPaste();

        Integer getMinViews();

        Integer getMaxViews();

        Integer getActivePastes();

        Integer getExpiredPastes();
}