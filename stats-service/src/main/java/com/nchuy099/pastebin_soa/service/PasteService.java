package com.nchuy099.pastebin_soa.service;

import com.nchuy099.pastebin_soa.repository.projection.MonthlyStatsProjection;

import java.time.YearMonth;

public interface PasteService {
    MonthlyStatsProjection getMonthlyStats(YearMonth yearMonth) throws Exception;
}
